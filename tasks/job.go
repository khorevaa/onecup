package tasks

type JobConfig struct {
	Steps []StepConfig `config:"steps,required" json:"steps"`
	Need  []string     `config:"need" json:"need"`
	If    string       `config:"if" json:"if"`
}
