package v1

type EnterpriseConfig struct {
	Dir          string `rawConfig:"dir,required" json:"dir"`
	FileTemplate string `rawConfig:"file-template" json:"file-template"`
}
