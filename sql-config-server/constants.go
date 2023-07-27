package sql_config_server

type Dialect uint

const (
	DialectPostgres Dialect = iota
	DialectMysql
	DialectSqlite
)
