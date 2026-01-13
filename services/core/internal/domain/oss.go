package domain

type UploadToken struct {
	Region   string `json:"region"`
	Endpoint string `json:"endpoint"`
	Bucket   string `json:"bucket"`
	Host     string `json:"host"`
	Dir      string `json:"dir"`
	Creds    Creds  `json:"credentials"`
}

type Creds struct {
	AccessKeyID     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	SecurityToken   string `json:"securityToken"`
	Expiration      string `json:"expiration"`
}
