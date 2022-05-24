package config

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
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
	IsDefaultsLocal bool
	// DefaultsFilePath string
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

const (
	defaultsFileName   = "xls-reports-defaults.yml"
	defaultsRemotePath = "https://gsales-lambdas.s3.ca-central-1.amazonaws.com/public/xls-reports-defaults.yml"
)

var (
	defs             = &defaults{}
	defaultsFilePath string
)

// ========================== Public Methods =============================== //

// Init method
func (c *Config) Init() (err error) {

	if err = c.setDefaults(); err != nil {
		return err
	}

	// I want the environment vars to be the final say, but we need them for the SSM Params
	// hence calling it twice
	if err = c.setEnvVars(); err != nil {
		return err
	}
	if err = c.setSSMParams(); err != nil {
		return err
	}

	if err = c.setEnvVars(); err != nil {
		return err
	}

	c.setDBConnectURL()
	c.setFinal()

	return err
}

// GetStageEnv method
func (c *Config) GetStageEnv() StageEnvironment {
	return c.Stage
}

// SetStageEnv method
func (c *Config) SetStageEnv(env string) (err error) {
	defs.Stage = env
	return c.validateStage()
}

// GetMongoConnectURL method
func (c *Config) GetMongoConnectURL() string {
	return c.DbConnectURL
}

// GetDbName method
func (c *Config) GetDbName() string {
	return c.config.DbName
}

// ========================== Private Methods =============================== //

// this must be called first in c.Load
func (c *Config) setDefaults() (err error) {

	var file []byte
	if c.IsDefaultsLocal == true { // DefaultsRemote is explicitly set to true

		dir, _ := os.Getwd()
		defaultsFilePath = path.Join(dir, defaultsFileName)
		if _, err = os.Stat(defaultsFilePath); os.IsNotExist(err) {
			return err
		}

		file, err = ioutil.ReadFile(defaultsFilePath)
		if err != nil {
			return err
		}

	} else { // using remote file path
		res, err := http.Get(defaultsRemotePath)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		file, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}
	}

	err = yaml.Unmarshal([]byte(file), &defs)
	if err != nil {
		return err
	}

	if err = c.validateStage(); err != nil {
		return err
	}

	return err
}

// validateStage method validates requested Stage exists
func (c *Config) validateStage() (err error) {

	validEnv := true

	switch defs.Stage {
	case "dev":
	case "development":
		c.Stage = DevEnv
	case "stage":
		c.Stage = StageEnv
	case "test":
		c.Stage = TestEnv
	case "prod":
		c.Stage = ProdEnv
	case "production":
		c.Stage = ProdEnv
	default:
		validEnv = false
	}

	if !validEnv {
		return errors.New(fmt.Sprintf("Invalid StageEnvironment requested: %s", defs.Stage))
	}

	return err
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
		Region: aws.String(defs.AwsRegion),
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
	if paramLen == 0 { // if no parameters returned then no sense in continueing
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

	if defs.DbUser != "" && defs.DbPassword != "" {
		userPass = fmt.Sprintf("%s:%s@", defs.DbUser, defs.DbPassword)
	}

	if userPass != "" {
		authSource = "?authSource=admin"
	}

	c.DbConnectURL = fmt.Sprintf("mongodb://%s%s/%s", userPass, defs.DbHost, authSource)

	return c
}

// Copies required fields from the defaults to the config struct
func (c *Config) setFinal() {
	c.AwsRegion = defs.AwsRegion
	c.CognitoClientID = defs.CognitoClientID
	c.DbName = defs.DbName
	c.S3Bucket = defs.S3Bucket
}
