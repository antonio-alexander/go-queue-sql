package sql

// environmental variables
const ()

// default configuration items
const (
	DefaultDatabaseHost  string = "localhost"
	DefaultDatabasePort  string = "3306"
	DefaultDatabaseName  string = "go_queue_sql"
	DefaultDatabaseTable string = "queue"
	DefaultDriverName    string = "mysql"
	DefaultUsername      string = "root"
	DefaultPassword      string = "mysql"
	DefaultParseTime     bool   = true
)

// Configuration provides the different items we can use to
// configure how we connect to the database
type Configuration struct {
	Hostname   string `json:"hostname"`
	Port       string `json:"port"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Database   string `json:"database"`
	Table      string `json:"table"`
	ParseTime  bool   `json:"parse_time"`
	DriverName string `json:"driver_name"`
}

func (c *Configuration) Validate() error {
	return nil
}

func (c *Configuration) Default() {
	c.Hostname = DefaultDatabaseHost
	c.Port = DefaultDatabasePort
	c.Database = DefaultDatabaseName
	c.Table = DefaultDatabaseTable
	c.DriverName = DefaultDriverName
	c.Username = DefaultUsername
	c.Password = DefaultPassword
	c.ParseTime = DefaultParseTime
}

//ConfigFromEnv can be used to generate a configuration pointer
// from a list of environments, it'll set the default configuraton
// as well

func (c *Configuration) FromEnvs(envs map[string]string) {
	if hostname, ok := envs["HOSTNAME"]; ok {
		c.Hostname = hostname
	}
	if port, ok := envs["PORT"]; ok {
		c.Port = port
	}
	if username, ok := envs["USERNAME"]; ok {
		c.Username = username
	}
	if password, ok := envs["PASSWORD"]; ok {
		c.Password = password
	}
	if database, ok := envs["DATABASE"]; ok {
		c.Database = database
	}
	if driver, ok := envs["DRIVER"]; ok {
		c.DriverName = driver
	}
}
