package sql_config_server

type Option func(*sqlConfigServer)

func WithTableCreation(dialect Dialect) Option {
	return func(server *sqlConfigServer) {
		server.dialect = &dialect
	}
}

func WithNamespace(ns string) Option {
	return func(server *sqlConfigServer) {
		server.namespace = &ns
	}
}
