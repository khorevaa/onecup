package tasks

type StepAction func(ctx Context, target Target, out map[string]interface{})

type Step struct {
	ID     string
	Name   string
	Params map[string]string
	Out    map[string]ParamConfig
	Action StepAction
}

type ParamConfig struct {
	Type string `json:"type" yaml:"type"`
}

type CacheConfig struct {
	Path string
}

type StepConfig struct {
	ID    string                 `json:"id,omitempty"`
	Name  string                 `config:"name,required" json:"name,omitempty"`
	With  map[string]string      `json:"with,omitempty"`
	Uses  string                 `config:"uses,required" json:"uses"`
	If    string                 `config:"if" json:"if"`
	Out   map[string]ParamConfig `json:"out"`
	Cache *CacheConfig           `config:"cache" json:"path"`
}
