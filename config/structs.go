package config

// defaults struct
type defaults struct {
	AwsRegion       string `yaml:"AwsRegion"`
	CognitoClientID string `yaml:"CognitoClientID"`
	DbHost          string `yaml:"DbHost"`
	DbName          string `yaml:"DbName"`
	DbPassword      string `yaml:"DbPassword"`
	DbUser          string `yaml:"DbUser"`
	S3Bucket        string `yaml:"S3Bucket"`
	SsmPath         string `yaml:"SsmPath"`
	Stage           string `yaml:"Stage"`
}

type config struct {
	AwsRegion       string
	CognitoClientID string
	DbConnectURL    string
	DbName          string
	S3Bucket        string
	Stage           StageEnvironment
}
