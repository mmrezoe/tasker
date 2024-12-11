package runner

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/mmrezoe/tasks/config"
)

func Run(c config.Config) error {

	if len(c.Vars) > 0 {
		for k, v := range c.Vars {
			c.Workflow.Vars[k] = v
		}
	}

	for k, v := range c.Workflow.Vars {
		if len(v) <= 0 {
			return fmt.Errorf("add var: %s", k)
		}
	}

	var wg sync.WaitGroup
	var mu sync.Mutex // Mutex to ensure thread-safe access to completed tasks
	errCh := make(chan error, len(c.Workflow.Tasks))

	var stdout, stderr *os.File
	var err error

	if len(c.OutFile) > 0 {
		stdout, err = os.OpenFile(c.OutFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("Error opening stdout file: %v", err)
		}
		defer stdout.Close()
	} else {
		stdout = os.Stdout
	}

	if len(c.ErrFile) > 0 {
		stderr, err = os.OpenFile(c.ErrFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("Error opening stderr file: %v", err)
		}
		defer stderr.Close()
	} else {
		stderr = os.Stderr
	}

	completedTasks := make(map[string]bool)

	for task := range c.Workflow.Tasks {
		wg.Add(1)

		if c.Workflow.Tasks[task].Concurrent {
			go ExecuteTask(c, task, stdout, stderr, &wg, completedTasks, &mu, errCh)
		} else {
			ExecuteTask(c, task, stdout, stderr, &wg, completedTasks, &mu, errCh)
		}
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for e := range errCh {
		if e != nil {
			return e
		}
	}

	return nil
}

func ExecuteTask(c config.Config, task int, stdout *os.File, stderr *os.File, wg *sync.WaitGroup, completedTasks map[string]bool, mu *sync.Mutex, errCh chan error) {
	defer wg.Done()

	for _, prereq := range c.Workflow.Tasks[task].Prerequisites {
		mu.Lock()
		if !completedTasks[prereq] {
			time.Sleep(time.Second * 1)
			mu.Unlock()

			config.Debug(fmt.Sprintf("Task %s waiting for prerequisite %s to complete.", c.Workflow.Tasks[task].Name, prereq), c.Debug)
			wg.Add(1)
			go ExecuteTask(c, task, stdout, stderr, wg, completedTasks, mu, errCh) // Re-execute this task later
			return
		}
		mu.Unlock()
	}

	config.Debug("Starting task: "+c.Workflow.Tasks[task].Name, c.Debug)
	for _, command := range c.Workflow.Tasks[task].Commands {
		command = ReplacePlaceholders(command, c.Workflow.Vars)

		cmd := exec.Command("sh", "-c", command)
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		err := cmd.Run()
		if err != nil {
			errCh <- fmt.Errorf("[Error] %s: %v\n", c.Workflow.Tasks[task].Name, err)
			return
		}
	}

	mu.Lock()
	completedTasks[c.Workflow.Tasks[task].Name] = true
	mu.Unlock()

	config.Debug("Completed task: "+c.Workflow.Tasks[task].Name, c.Debug)
}

func ReplacePlaceholders(command string, variables map[string]string) string {

	for key, value := range variables {
		placeholder := fmt.Sprintf("{{%s}}", key)
		command = strings.ReplaceAll(command, placeholder, value)
	}

	return command
}
