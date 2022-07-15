package config

import "time"

// defaults struct
type defaults struct {
	AwsRegion  string `yaml:"AwsRegion"`
	DbHost     string `yaml:"DbHost"`
	DbName     string `yaml:"DbName"`
	DbPassword string `yaml:"DbPassword"`
	DbUser     string `yaml:"DbUser"`
	ExpireHrs  int    `yaml:"ExpireHrs"`
	S3Bucket   string `yaml:"S3Bucket"`
	SsmPath    string `yaml:"SsmPath"`
	Stage      string `yaml:"Stage"`
}

type config struct {
	AwsRegion     string
	DbConnectURL  string
	DbName        string
	S3Bucket      string
	Stage         StageEnvironment
	UrlExpireTime time.Duration
}
