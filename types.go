package sql

import "time"

const (
	DefaultDatabaseHost  string        = "localhost"
	DefaultDatabasePort  string        = "3306"
	DefaultDatabaseName  string        = "go_queue_sql"
	DefaultDatabaseTable string        = "queue"
	DefaultDriverName    string        = "mysql"
	DefaultUsername      string        = "root"
	DefaultPassword      string        = "mysql"
	DefaultCreateTable   bool          = false
	DefaultParseTime     bool          = true
	DefaultQueryTimeout  time.Duration = 30 * time.Second
)

//Configuration provides the different items we can use to
// configure how we connect to the database
type Configuration struct {
	Hostname     string        `json:"hostname"`
	Port         string        `json:"port"`
	Username     string        `json:"username"`
	Password     string        `json:"password"`
	Database     string        `json:"database"`
	Table        string        `json:"table"`
	ParseTime    bool          `json:"parse_time"`
	DriverName   string        `json:"driver_name"`
	CreateTable  bool          `json:"create_table"`
	QueryTimeout time.Duration `json:"query_timeout"`
}

type bytes []byte

type Owner interface {
	Initialize(config *Configuration) error
	Shutdown() error
}
