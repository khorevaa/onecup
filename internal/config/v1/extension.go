package v1

type ExtensionConfig struct {
	Dir          string `config:"dir,required" json:"dir"`
	FileTemplate string `config:"file-template" json:"file-template"`
}
