package config

var OPTS struct {
	LogLevel string `short: "logLevel" long:"logLevel"`
	DBConn string `short: "db" long:"databaseConnectionString"`
	DBDriver string `short: "driver" long: "databaseDriver"`
	MaxDatabaseConnections int `short: "maxDatabaseConnections", long: "MaxDatabaseConnections"`
	SqlCACertFile string `short: "sqlCACertFile long: "sqlCACertFile"`
}
