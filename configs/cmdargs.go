package configs

import (
	"flag"
	"log"
)

const (
	// environment declarations the app is running on
	EnvDev     = "development"
	EnvStaging = "staging"
	EnvProd    = "production"
)


// CmdArgsConfig has all configurations passed to the app as
// cmd arguments
type CmdArgsConfig struct {
	Environment string // 'development' <alias: dev>, 'staging', 'production' <alias: prod>
}

// ParseCMDArgs gets all flags passed to the application as arguments
// into a config struct
func ParseCMDArgs() CmdArgsConfig {
	var env string
	var config CmdArgsConfig

	flag.StringVar(&env, "env", EnvDev, "Running environment of application, dev, staging, prod")
	flag.Parse()

	log.Println(env)

	switch env {
	case "dev", "development":
		config.Environment = EnvDev
	case "staging":
		config.Environment = EnvStaging
	case "prod", "production":
		config.Environment = EnvProd
	}

	return config
}
