package config

type Config struct {
	Workflow Workflow
	Vars     map[string]string
	Debug    bool
	OutFile  string
	ErrFile  string
}

type Workflow struct {
	Name   string            `yaml:"name"`
	Author string            `yaml:"author"`
	Usage  string            `yaml:"usage"`
	Tasks  []Task            `yaml:"tasks"`
	Vars   map[string]string `yaml:"vars"`
}

type Task struct {
	Name          string   `yaml:"name"`
	Image         string   `yaml:"image"`
	Commands      []string `yaml:"commands"`
	Prerequisites []string `yaml:"prerequisites,omitempty"`
	Concurrent    bool     `yaml:"concurrent,omitempty"`
}
