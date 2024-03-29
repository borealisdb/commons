package postgresql

import (
	"context"
	"fmt"
	"github.com/borealisdb/commons/credentials"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"os"
	"strings"
)

type Postgresql interface {
	GetConnection(ctx context.Context, clusterName string, username string, options Options) (*sqlx.DB, error)
	Connect(dsn string) (*sqlx.DB, error)
	GetDSN(clusterName, password, username string, options Options) (string, error)
	GetCredentials(ctx context.Context, clusterName string, username string, options Options) (credentials.GetPostgresCredentialsResponse, error)
}

type Options struct {
	Database        string
	Port            string
	Host            string // If host is not specified it will work with clusterName
	Role            string // master or replica
	SSLRootCertPath string
	SSLMode         string
	SSLDownload     bool

	SetMaxIdleConns    int
	SetMaxOpenConns    int
	SetConnMaxLifetime int
}

type PG struct {
	CredentialsProvider credentials.Credentials
	options             Options
	clusterName         string
}

func (pg *PG) GetConnection(ctx context.Context, clusterName string, username string, options Options) (*sqlx.DB, error) {
	if err := pg.setDefaults(options, clusterName); err != nil {
		return &sqlx.DB{}, err
	}
	resp, err := pg.getCredentials(ctx, username)
	if err != nil {
		return &sqlx.DB{}, fmt.Errorf("could not GetCredentials: %v", err)
	}

	if pg.options.SSLDownload {
		if err := pg.downloadSSLRootCert(ctx); err != nil {
			return &sqlx.DB{}, fmt.Errorf("could not downloadSSLRootCert: %v", err)
		}
	}

	dsn := pg.getDSN(resp.Password, resp.Username)
	return pg.Connect(dsn)
}

func (pg *PG) Connect(dsn string) (*sqlx.DB, error) {
	conn, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return &sqlx.DB{}, err
	}
	conn.SetMaxIdleConns(1)
	conn.SetMaxOpenConns(1)
	conn.SetConnMaxLifetime(0)
	return conn, nil
}

func (pg *PG) getDSN(password string, username string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(
		"postgresql://%v:%v@%v:%v/%v?sslmode=%v",
		username,
		password,
		pg.options.Host,
		pg.options.Port,
		pg.options.Database,
		pg.options.SSLMode,
	))
	if pg.options.SSLRootCertPath != "" {
		sb.WriteString(fmt.Sprintf("&sslrootcert=%v", pg.options.SSLRootCertPath))
	}

	return sb.String()
}

func (pg *PG) GetDSN(clusterName, password, username string, options Options) (string, error) {
	if err := pg.setDefaults(options, clusterName); err != nil {
		return "", fmt.Errorf("could not setDefaults: %v", err)
	}

	return pg.getDSN(password, username), nil
}

func (pg *PG) GetCredentials(ctx context.Context, clusterName string, username string, options Options) (credentials.GetPostgresCredentialsResponse, error) {
	if err := pg.setDefaults(options, clusterName); err != nil {
		return credentials.GetPostgresCredentialsResponse{}, fmt.Errorf("could not setDefaults: %v", err)
	}
	return pg.getCredentials(ctx, username)
}

func (pg *PG) getCredentials(ctx context.Context, username string) (credentials.GetPostgresCredentialsResponse, error) {
	postgresCredentials, err := pg.CredentialsProvider.GetPostgresCredentials(
		ctx,
		pg.clusterName,
		username,
		credentials.Options{},
	)
	if err != nil {
		return credentials.GetPostgresCredentialsResponse{}, fmt.Errorf("could not GetPostgresCredentials: %v", err)
	}

	return postgresCredentials, nil
}

func (pg *PG) downloadSSLRootCert(ctx context.Context) error {
	cert, err := pg.CredentialsProvider.GetPostgresSSLRootCert(ctx, pg.clusterName, credentials.Options{})
	if err != nil {
		return fmt.Errorf("could not GetPostgresSSLRootCert: %v", err)
	}

	return os.WriteFile(pg.options.SSLRootCertPath, cert.RootCertBytes, 0777)
}

func (pg *PG) setDefaults(options Options, clusterName string) error {
	pg.clusterName = clusterName
	pg.options = options
	resp, err := pg.CredentialsProvider.GetClusterEndpoint(context.Background(), pg.clusterName, pg.options.Role)
	if pg.options.Host == "" {
		if err != nil {
			return fmt.Errorf("could not GetClusterEndpoint: %v", err)
		}
		pg.options.Host = resp.Hostname
	}
	if pg.options.Port == "" {
		pg.options.Port = resp.Port
	}
	if pg.options.Database == "" {
		pg.options.Database = "postgres"
	}

	if pg.options.SSLRootCertPath == "" && pg.options.SSLMode == "" {
		pg.options.SSLMode = "disable"
	} else if pg.options.SSLRootCertPath != "" && pg.options.SSLMode == "" {
		pg.options.SSLMode = "verify-ca"
	}

	return nil
}
