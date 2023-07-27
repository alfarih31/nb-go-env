package sql_config_server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	env "github.com/alfarih31/nb-go-env"
	"github.com/alfarih31/nb-go-env/internal"
	"log"
)

type sqlConfigServer struct {
	tableName string
	namespace *string
	db        *sql.DB
	dialect   *Dialect
}

type SqlConfigServer interface {
	env.ConfigServer
}

func (s *sqlConfigServer) Get(key string) (string, bool) {
	var r *sql.Row
	if s.namespace != nil {
		r = s.db.QueryRow(fmt.Sprintf(`SELECT value FROM "%s" WHERE key = '%s' AND namespace = '%s' ORDER BY id DESC LIMIT 1`, s.tableName, key, *s.namespace))
	} else {
		r = s.db.QueryRow(fmt.Sprintf(`SELECT value FROM "%s" WHERE key = '%s' ORDER BY id DESC LIMIT 1`, s.tableName, key))
	}

	var v string
	if err := r.Scan(&v); err != nil {
		log.Println(err)
		return "", false
	}

	return v, !internal.HasZeroValue(v)
}

func (s *sqlConfigServer) Dump() (string, error) {
	var (
		rs  *sql.Rows
		err error
	)
	if s.namespace != nil {
		rs, err = s.db.Query(fmt.Sprintf(`SELECT key, value FROM "%s" WHERE namespace = '%s' ORDER BY id DESC`, s.tableName, *s.namespace))
	} else {
		rs, err = s.db.Query(fmt.Sprintf(`SELECT key, value FROM "%s" ORDER BY id DESC`, s.tableName))
	}

	if err != nil {
		return "", err
	}

	configs := map[string]string{}
	for rs.Next() {
		var k, v string
		if err = rs.Scan(&k, &v); err != nil {
			return "", err
		}
		configs[k] = v
	}

	if err = rs.Close(); err != nil {
		return "", err
	}

	j, e := json.Marshal(configs)

	return string(j), e
}

func (s *sqlConfigServer) executeTableCreation() error {
	if s.dialect != nil {
		var createStmt string
		switch *s.dialect {
		case DialectPostgres:
			createStmt = fmt.Sprintf(`
  CREATE TABLE IF NOT EXISTS "%s" (
  id SERIAL PRIMARY KEY,
  namespace CHARACTER VARYING(511),
  key VARCHAR(511) NOT NULL,
  value TEXT NOT NULL
  );`, s.tableName)
		case DialectMysql:
			createStmt = fmt.Sprintf(`
  CREATE TABLE IF NOT EXISTS "%s" (
  id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
  namespace VARCHAR(511),
  key VARCHAR(511) NOT NULL,
  value TEXT NOT NULL
  );`, s.tableName)
		case DialectSqlite:
			createStmt = fmt.Sprintf(`
  CREATE TABLE IF NOT EXISTS "%s" (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  namespace VARCHAR(511),
  key VARCHAR(511) NOT NULL,
  value TEXT NOT NULL
  );`, s.tableName)
		default:
			return errors.New("unknown dialect")
		}

		_, err := s.db.Exec(createStmt)
		if err != nil {
			return err
		}

		return nil
	}

	return nil
}

func (s *sqlConfigServer) init() error {
	if err := s.executeTableCreation(); err != nil {
		return err
	}

	return nil
}

func NewSqlConfigServer(db *sql.DB, tableName string, opts ...Option) (SqlConfigServer, error) {
	cs := &sqlConfigServer{
		db:        db,
		tableName: tableName,
	}

	for _, opt := range opts {
		opt(cs)
	}

	if err := cs.init(); err != nil {
		return nil, err
	}

	return cs, nil
}
