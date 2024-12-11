<h1 align="center">
  Tasks
</h1>

A command line automation tool to speed up tasks. Perfect for automating and putting together tools for bug bounty hackers. Inspired by [RikunjSindhwad/Task-Ninja](https://github.com/RikunjSindhwad/Task-Ninja)

## notes

- If a variable is mandatory and its value must be entered manually, it is sufficient to set its value to nothing in the yaml file.

## example usages

```bash
tasks -w test-workflow.yaml -v age=20 -d -o stdout.txt -e stderr.txt
```

or print usage if workflow

```bash
tasks -w test-workflow.yaml -u
```

## Installation

- Go Install
  ```bash
  GO111MODULE=on
  go install github.com/mmrezoe/tasks@latest
  ```
- Build

  ```bash
  # Clone the repository
  git clone https://github.com/mmrezoe/tasks.git

  # Navigate to the Tasks directory
  cd tasks

  # Build Tasks
  go build
  ```
