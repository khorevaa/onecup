package tasks

type Config struct {
	Context ContextConfig `yaml:"context"`
	Targets []TargetConfig
	Jobs    map[string]JobConfig
}
