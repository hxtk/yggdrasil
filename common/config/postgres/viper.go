package postgres

import (
	"database/sql"
	"fmt"
	"net/url"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/viper"
)

func FromViper(v *viper.Viper) (*sql.DB, error) {
	postgresURL := &url.URL{
		Scheme: "postgres",
	}
	if user := v.GetString("postgres.user"); user != "" {
		if password := v.GetString("postgres.password"); password != "" {
			postgresURL.User = url.UserPassword(user, password)
		} else {
			postgresURL.User = url.User(user)
		}
	}

	if host := v.GetString("postgres.host"); host != "" {
		if port := v.GetString("postgres.port"); port != "" {
			postgresURL.Host = fmt.Sprintf("%s:%s", host, port)
		} else {
			postgresURL.Host = host
		}
	}

	if database := v.GetString("postgres.database"); database != "" {
		postgresURL.Path = database
	}

	values := make(url.Values)
	if sslmode := v.GetString("postgres.sslmode"); sslmode != "" {
		values.Set("sslmode", sslmode)
	}

	if sslcert := v.GetString("postgres.sslcert"); sslcert != "" {
		values.Set("sslcert", sslcert)
	}

	if sslkey := v.GetString("postgres.sslkey"); sslkey != "" {
		values.Set("sslkey", sslkey)
	}

	if sslrootcert := v.GetString("postgres.sslrootcert"); sslrootcert != "" {
		values.Set("sslrootcert", sslrootcert)
	}

	postgresURL.RawQuery = values.Encode()
	dsn := postgresURL.String()
	if v.IsSet("postgres.schema") {
		m, err := migrate.New(v.GetString("postgres.schema"), dsn)
		if err != nil {
			return nil, err
		}
		if err := m.Up(); err != nil {
			return nil, err
		}
	}

	return sql.Open("postgres", dsn)
}
