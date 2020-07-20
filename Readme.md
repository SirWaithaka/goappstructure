# Golang Apps Structure

Question: How do you structure your Go applications in a manner that doesn't impose too much verbosity for a small
budding idea, but still scales well for an enterprise scale solution?

I have found this structure to be simple, clean and scalable.

## The Structure

```
.
+-- app
|   +-- entities
|   +-- errorz
|   +-- routing
|   +-- params
+-- configs
+-- database
+-- migrations
+-- registry
+-- config.yml
+-- main.go
```

### configs - dir
This directory has any configuration objects that represent any external configurations files, e.g. yaml/toml and
dynamic configurations passed to the applications as cmd arguments.

Example of a config file using yaml structure

```yaml
application:
    port: 2801

database:
    user: "sirwaithaka"
    dbname: "test"
    password: "sirwaithaka"
    port: 5432
    driver: "postgresql"
    host: "localhost"
```

If we have this example configuration file, we can map it to a struct as so inside the configs dir.

```go
// configs/yaml.go

type YamlConfig struct {
	Application struct {
        Port int `yaml:"port"`
    } `yaml:"application"`
    Database struct {
        User string `yaml:"user"`
        DbName string `yaml:"dbname"`
        Password string `yaml:"password"`
        Port int `yaml:"port"`
        Driver string `yaml:"driver"`
        Host string  `yaml:"host"`
    } `yaml:"database"`
}
```

Example of a configuration parsed from cmd arguments for an application that takes as an argument the environment of the
application, then uses that to determine which `yaml` configuration file to run. This assuming we are using different
`config.yml` setups e.g. `config.dev.yml` for the `dev` environment, `config.staging.yml` for the `staging` environment
and `config.yml` for the `prod` environment.

Example struct mapping

```go
// configs/cmdargs.go


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
    
    // by default we run the app in dev mode
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
```

### database - dir
This directory has the wrapper struct for a database connection, and the logic to create a connection to the database. In
this case lets assume we are using `gorm` as our data access layer.

```
+-- database
|   +-- database.go
|   +-- gorm.go
```

```go
// database/gorm.go

// NewGormConnection creates a connection to the database and returns
// a gorm.DB pointer object or an error
func NewGormConnection(dbString string) (*gorm.DB, error) {

	conn, err := gorm.Open("postgres", dbString)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
```

First we create a connection to the database using `gorm`.

```go
// database/database.go

// Database is our application level wrapper of connection pool
// to the  database.
type Database struct {
	*gorm.DB // gorm connection object
}

// NewDatabase creates a new Database wrapper
func NewDatabase(connectionString string) (*Database, error) {
	conn, err := NewGormConnection(connectionString)
	if err != nil {
		return nil, err
	}

	conn.LogMode(true)

	return &Database{conn}, nil
}
```

We initialize our data access layer with our own database wrapper struct, which we can add methods to if needed. We can
also extend our database wrapper to have connections to other forms of databases e.g. `redis` cache.

### migrations - dir
This directory has any sql migrations and possible seed data that we can use to populate a working database for the
application. You can use whichever tool suites your poison.

1. [Goose](https://github.com/pressly/goose/cmd/goose   )
2. [Soda](https://github.com/gobuffalo/pop)
3. You can find other tools online.

### app - dir
This directory contains all your business logic following a clean architecture pattern (sort of).

```
+-- app
|   +-- entities
|   +-- errorz
|   +-- params
|   +-- routing
```

#### app/entities - dir
Entities contain `struct`s that define a mapping of our sql tables and any other models that the application domain
has.

Example 

```go
// app/entities/user.go

type User struct {
    ID uuid.UUID `gorm:"column:user_id"`
    FirstName string `gorm:"column:first_name"`
    LastName string `gorm:"column:last_name"`
}
```

#### app/errorz - dir
The naming of the dir with a 'z' is upon user preference. Choose your own poison. However the choice here is not to
collide with std lib's `errors` package.

This directory has global error handling logic. You can categorize errors depending on where in the application layer
they occur. Example errors at the sql or data access layer, errors at the business logic layer or errors at the http
handler layer / (`views`).

#### app/params - dir
This directory defines the parameters that the application expects as inputs. Example parameters from `POST` or `PUT`
requests to the application. These same parameters will be used as inputs to the business logic layer for further
business processing.

Example
1. Let's say we want to create a user, this will be a `POST` request with some `form` values

```go
// app/params/users.go

type CreateUser struct {
    FirstName string `json:"firstName" form:"firstName" validate:"required"`
    LastName string `json:"lastName" form:"lastName" validate:"required"`
    Age int `json:"age" form:"age"`
}
```

These fields can be parsed and those with the tag `required` can be validated and then the object can be used as argument
in the business logic to create the user.

#### app/routing -  dir
This directory defines the `http` routes used in the api or application. Inside the dir you can have any middleware logic
for your http requests.