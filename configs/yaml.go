package configs

import (
	"os"

	"gopkg.in/yaml.v2"
)

const (
	// declaration of yaml file paths for different environments
	yamlDevPath     = "config.dev.yml"
	yamlStagingPath = "config.staging.yml"
	yamlProdPath    = "config.yml"
)

var envYamlMapping = map[string]string{
	EnvDev:     yamlDevPath,
	EnvStaging: yamlStagingPath,
	EnvProd:    yamlProdPath,
}

// YamlConfig is a mapping struct with same properties
// existing in the configuration yaml file.
type YamlConfig struct {
	Application struct {
		Port int `yaml:"port"`
	} `yaml:"application"`

	Database struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`

		Name   string `yaml:"name"`
		Host   string `yaml:"host"`
		Port   int    `yaml:"port"`
		Client string `yaml:"client"`

		SSLCertPath  string `yaml:"ssl_cert_path"`
		SSLKeyPath   string `yaml:"ssl_key_path"`
		RootCertPath string `yaml:"root_cert_path"`
	} `yaml:"database"`
}

// ReadYaml uses commandline arguments flags checks the env
// and reads the yaml file mapping to the env
func ReadYaml(cmdArgs CmdArgsConfig) (YamlConfig, error) {
	yamlPath := envYamlMapping[cmdArgs.Environment]

	// read yaml config file
	f, err := os.Open(yamlPath)
	if err != nil {
		return YamlConfig{}, err
	}
	defer f.Close()

	var config YamlConfig
	err = yaml.NewDecoder(f).Decode(&config)
	if err != nil {
		return YamlConfig{}, err
	}

	return config, nil
}
