package postgresql

import (
	"fmt"
	"github.com/borealisdb/commons/constants"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"strings"
)

type PostgresqlV2 interface {
	GetConnection(args Args) (*sqlx.DB, error)
	Connect(dsn string) (*sqlx.DB, error)
	GetDSN(args Args) (string, error)
}

type Args struct {
	Username string
	Password string
	Database string
	Port     string
	Host     string

	SSLRootCertPath string
	SSLMode         string
}

type V2 struct{}

func (pg *V2) GetConnection(args Args) (*sqlx.DB, error) {
	argsWithDefaults := pg.setDefaults(args)
	dsn := pg.getDSN(argsWithDefaults)
	return pg.Connect(dsn)
}

func (pg *V2) Connect(dsn string) (*sqlx.DB, error) {
	conn, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return &sqlx.DB{}, err
	}
	conn.SetMaxIdleConns(1)
	conn.SetMaxOpenConns(1)
	conn.SetConnMaxLifetime(0)
	return conn, nil
}

func (pg *V2) GetDSN(args Args) (string, error) {
	argsWithDefaults := pg.setDefaults(args)
	return pg.getDSN(argsWithDefaults), nil
}

func (pg *V2) getDSN(args Args) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(
		"postgresql://%v:%v@%v:%v/%v?sslmode=%v",
		args.Username,
		args.Password,
		args.Host,
		args.Port,
		args.Database,
		args.SSLMode,
	))
	if args.SSLRootCertPath != "" {
		sb.WriteString(fmt.Sprintf("&sslrootcert=%v", args.SSLRootCertPath))
	}

	return sb.String()
}

func (pg *V2) setDefaults(args Args) Args {
	newArgs := args
	if args.Host == "" {
		newArgs.Host = "localhost"
	}
	if args.Port == "" {
		newArgs.Port = constants.PostgresDefaultPort
	}
	if args.Database == "" {
		newArgs.Database = "postgres"
	}

	if args.SSLMode == "" {
		newArgs.SSLMode = "disable"
	} else if args.SSLRootCertPath != "" && args.SSLMode == "" {
		newArgs.SSLMode = "verify-ca"
	}

	return newArgs
}
