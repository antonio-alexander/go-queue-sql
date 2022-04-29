package sql

const (
	defaultDatabaseHost  string = "localhost"
	defaultDatabasePort  string = "3306"
	defaultDatabaseName  string = "go_queue_sql"
	defaultDatabaseTable string = "queue"
	defaultDriverName    string = "mysql"
	defaultUsername      string = "root"
	defaultPassword      string = "mysql"
)

//Configuration provides the different items we can use to
// configure how we connect to the database
type Configuration struct {
	Hostname   string `json:"hostname"`   //hostame to user to access the database
	Port       string `json:"port"`       //port to connect to
	Username   string `json:"username"`   //username to authenticate with
	Password   string `json:"password"`   //password to authenticate with
	Database   string `json:"database"`   //database to connect to
	Table      string `json:"table"`      //
	ParseTime  bool   `json:"parse_time"` //whether or not to parse time
	DriverName string `json:"driver_name"`
}

//ConfigFromEnv can be used to generate a configuration pointer
// from a list of environments, it'll set the default configuraton
// as well
func ConfigFromEnv(envs map[string]string) *Configuration {
	c := &Configuration{
		Hostname:   defaultDatabaseHost,
		Port:       defaultDatabasePort,
		Database:   defaultDatabaseName,
		Table:      defaultDatabaseTable,
		DriverName: defaultDriverName,
		Username:   defaultUsername,
		Password:   defaultPassword,
		ParseTime:  false,
	}
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
	return c
}
