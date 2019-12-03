package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"gopkg.in/yaml.v2"
)

// Config struct
type Config struct {
	config
	DefaultsFilePath string
}

// StageEnvironment string
type StageEnvironment string

// StageEnvironment type constants
const (
	DevEnv   StageEnvironment = "dev"
	StageEnv StageEnvironment = "stage"
	TestEnv  StageEnvironment = "test"
	ProdEnv  StageEnvironment = "prod"
)

const defaultFileName = "defaults.yml"

var (
	defs = &defaults{}
)

// Load method
func (c *Config) Load() (err error) {

	err = c.setDefaults()
	if err != nil {
		return err
	}
	err = c.setEnvVars()
	if err != nil {
		return err
	}
	err = c.setSSMParams()
	if err != nil {
		return err
	}
	err = c.setFinal()
	if err != nil {
		return err
	}

	c.setDBConnectURL()
	return err
}

// GetStageEnv method
func (c *Config) GetStageEnv() StageEnvironment {
	return c.Stage
}

// GetMongoConnectURL method
func (c *Config) GetMongoConnectURL() string {
	return c.DBConnectURL
}

// this must be called first in c.Load
func (c *Config) setDefaults() (err error) {

	if c.DefaultsFilePath == "" {
		dir, _ := os.Getwd()
		c.DefaultsFilePath = path.Join(dir, defaultFileName)
	}

	file, err := ioutil.ReadFile(c.DefaultsFilePath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal([]byte(file), &defs)
	if err != nil {
		return err
	}
	err = c.validateStage()

	return err
}

// validateStage method validates requested Stage exists
func (c *Config) validateStage() error {

	validEnv := false

	switch defs.Stage {
	case "dev":
		c.Stage = DevEnv
		validEnv = true
	case "stage":
		c.Stage = StageEnv
		validEnv = true
	case "test":
		c.Stage = TestEnv
		validEnv = true
	case "prod":
		c.Stage = ProdEnv
		validEnv = true
	case "production":
		c.Stage = ProdEnv
		validEnv = true
	}

	if !validEnv {
		// return errors.New(fmt.Sprintf("Invalid Stage type: %s", c.Stage))
		return fmt.Errorf("Invalid Stage type: %s", defs.Stage)
	}

	return nil
}

// sets any environment variables that match the default struct fields
func (c *Config) setEnvVars() (err error) {

	vals := reflect.Indirect(reflect.ValueOf(defs))
	for i := 0; i < vals.NumField(); i++ {
		nm := vals.Type().Field(i).Name
		if e := os.Getenv(nm); e != "" {
			vals.Field(i).SetString(e)
		}
		// If field is Stage, validate and return error if required
		if nm == "Stage" {
			err = c.validateStage()
			if err != nil {
				return err
			}
		}
	}

	return err
}

func (c *Config) setSSMParams() (err error) {

	s := []string{"", string(c.GetStageEnv()), defs.SsmPath}
	paramPath := aws.String(strings.Join(s, "/"))

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(defs.AWSRegion),
	})
	if err != nil {
		return err
	}

	svc := ssm.New(sess)
	res, err := svc.GetParametersByPath(&ssm.GetParametersByPathInput{
		Path:           paramPath,
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return err
	}

	paramLen := len(res.Parameters)
	if paramLen == 0 {
		// err = fmt.Errorf("Error fetching ssm params, total number found: %d", paramLen)
		return nil
	}

	// Get struct keys so we can test before attempting to set
	t := reflect.ValueOf(defs).Elem()
	for _, r := range res.Parameters {
		paramName := strings.Split(*r.Name, "/")[3]
		structKey := t.FieldByName(paramName)
		if structKey.IsValid() {
			structKey.Set(reflect.ValueOf(*r.Value))
		}
	}
	return err
}

// Build a url used in mgo.Dial as described in: https://godoc.org/gopkg.in/mgo.v2#Dial
func (c *Config) setDBConnectURL() *Config {

	var userPass, authSource string

	if defs.DBUser != "" && defs.DBPassword != "" {
		userPass = defs.DBUser + ":" + defs.DBPassword + "@"
	}

	if userPass != "" {
		authSource = "?authSource=admin"
	}

	c.DBConnectURL = "mongodb://" + userPass + defs.DBHost + "/" + authSource

	return c
}

// Copies required fields from the defaults to the Config struct
func (c *Config) setFinal() (err error) {

	c.AWSRegion = defs.AWSRegion
	c.S3Bucket = defs.S3Bucket
	c.DBName = defs.DBName
	err = c.validateStage()

	return err
}
