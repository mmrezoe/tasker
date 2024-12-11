package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

func ParseArgs() Config {
	version := "0.0.1"
	var workflowFile string
	flag.StringVar(&workflowFile, "w", "", "Path to the workflow YAML file")
	var varStr string
	flag.StringVar(&varStr, "v", "", "Comma-separated list of vars (key=value,key2=value2)")
	var outFile string
	flag.StringVar(&outFile, "o", "", "Stdout file name(If there is no value, it should print the output.)")
	var errFile string
	flag.StringVar(&errFile, "e", "", "Stderr file name(If there is no value, it should print the output.)")
	var usageWF bool
	flag.BoolVar(&usageWF, "u", false, "usage of workflow")
	var debug bool
	flag.BoolVar(&debug, "d", false, "debug mode")
	flag.Usage = func() {
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "version: %s\n", version)
	}
	flag.Parse()

	if len(workflowFile) <= 0 {
		fmt.Println("Please specify a workflow file using the -w flag.")
		os.Exit(1)
	}

	workflow, err := ReadYamlFile(workflowFile)
	if err != nil {
		log.Fatalf("Error reading workflow file: %v\n", err)
	}

	if usageWF {
		fmt.Printf("Usage:\n- %s\n", workflow.Usage)
		if len(workflow.Vars) > 0 {
			fmt.Println("Default vars")

			for k, v := range workflow.Vars {
				fmt.Printf("- %s: %s\n", k, v)
			}
		}
		os.Exit(0)
	}

	variables := make(map[string]string)
	if varStr != "" {
		vars := strings.Split(varStr, ",")
		for _, v := range vars {
			kv := strings.Split(v, "=")
			if len(kv) == 2 {
				variables[kv[0]] = kv[1]
			}
		}
	}

	// return workflow, variables, debug
	return Config{
		Workflow: workflow,
		Vars:     variables,
		Debug:    debug,
		OutFile:  outFile,
		ErrFile:  errFile,
	}
}

// ReadYamlFile reads the workflow YAML file
func ReadYamlFile(filePath string) (Workflow, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Workflow{}, err
	}
	defer file.Close()

	var workflow Workflow
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&workflow)
	if err != nil {
		return Workflow{}, err
	}

	return workflow, nil
}
