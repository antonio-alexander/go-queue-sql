package sql

import "strconv"

//ConfigFromEnv can be used to generate a configuration pointer
// from a list of environments, it'll set the default configuraton
// as well
func ConfigFromEnv(envs map[string]string) *Configuration {
	c := &Configuration{
		Hostname:     DefaultDatabaseHost,
		Port:         DefaultDatabasePort,
		Database:     DefaultDatabaseName,
		Table:        DefaultDatabaseTable,
		DriverName:   DefaultDriverName,
		Username:     DefaultUsername,
		Password:     DefaultPassword,
		ParseTime:    DefaultParseTime,
		CreateTable:  DefaultCreateTable,
		QueryTimeout: DefaultQueryTimeout,
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
	if s, ok := envs["CREATE_TABLE"]; ok {
		if createTable, err := strconv.ParseBool(s); err != nil {
			c.CreateTable = createTable
		}
	}
	return c
}
