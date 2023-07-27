package main

import (
	"database/sql"
	"fmt"
	env2 "github.com/alfarih31/nb-go-env"
	sql_config_server "github.com/alfarih31/nb-go-env/sql-config-server"
	_ "github.com/lib/pq"
)

const (
	host      = "localhost"
	port      = 5432
	user      = "postgres"
	password  = "postgres"
	dbname    = "postgres"
	namespace = "example"
)

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	// create config server
	configService, err := sql_config_server.NewSqlConfigServer(db, "configs",
		sql_config_server.WithNamespace(namespace),
		sql_config_server.WithTableCreation(sql_config_server.DialectPostgres))

	CheckError(err)

	// seed the database
	insertStmt := fmt.Sprintf(`insert into "configs"("key", "value", "namespace") values('foo', 'ba2', '%s')`, namespace)
	_, e := db.Exec(insertStmt)
	CheckError(e)

	env, err := env2.LoadWithConfigServer(configService)
	CheckError(err)

	fmt.Printf("Config with key 'foo' -> '%s'\n", env.MustGetString("foo"))

	db.Close()
}
