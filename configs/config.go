package configs

import "fmt"

// Database
type Database struct {
	User     string

	Host   string
	Port   int
	DBName string
	Client string

	SSL      bool
	SSLMode  string
	SSLCert  string
	SSLKey   string
	RootCert string
}

// String returns a connection string as require by lib/pq library see
// https://godoc.org/github.com/lib/pq#hdr-Connection_String_Parameters
func (d Database) String() string {
	return fmt.Sprintf("%s://%s@%s:%d/%s?ssl=%v&sslmode=%s&sslrootcert=%s&sslkey=%s&sslcert=%s",
		d.Client, d.User, d.Host, d.Port, d.DBName, d.SSL, d.SSLMode, d.RootCert, d.SSLKey, d.SSLCert)
}

// Stringify builds a connection string as required by lib/pq driver see
// https://godoc.org/github.com/lib/pq#hdr-Connection_String_Parameters
func (d Database) Stringify() string {
	return fmt.Sprintf(""+
		"user=%v "+
		"dbname=%v "+
		"host=%v "+
		"port=%v "+
		"sslmode=%v "+
		"sslrootcert=%v "+
		"sslkey=%v "+
		"sslcert=%v "+
		"", d.User, d.DBName, d.Host, d.Port, d.SSLMode, d.RootCert, d.SSLKey, d.SSLCert)
}


// Config is application level configuration
// it is the config struct/object that is passed around
// into the business logic and instantiation of other
// middleware objects.
type Config struct {
	Port int

	DB Database
}

// GetConfig interprets the yaml configuration into an app config instance
func GetConfig(cfg YamlConfig) *Config {
	return &Config{
		Port: cfg.Application.Port,
		DB:   Database{
			Client:   cfg.Database.Client,
			User:     cfg.Database.User,
			Host:     cfg.Database.Host,
			Port:     cfg.Database.Port,
			DBName:   cfg.Database.Name,

			SSL:      true,
			SSLMode:  "require",
			SSLCert:  cfg.Database.SSLCertPath,
			SSLKey:   cfg.Database.SSLKeyPath,
			RootCert: cfg.Database.RootCertPath,
		},
	}
}
