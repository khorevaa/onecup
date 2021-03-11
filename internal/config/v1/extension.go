package v1

type ExtensionConfig struct {
	Dir          string `rawConfig:"dir,required" json:"dir"`
	FileTemplate string `rawConfig:"file-template" json:"file-template"`
}
