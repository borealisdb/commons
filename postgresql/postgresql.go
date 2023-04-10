package postgresql

import (
	"context"
	"fmt"
	"github.com/borealisdb/commons/credentials"
	"github.com/borealisdb/commons/environment"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"os"
	"strings"
)

type Postgresql interface {
	GetConnection(ctx context.Context, clusterName, username string, options Options) (*sqlx.DB, error)
	Connect(dsn string) (*sqlx.DB, error)
	GetPostgresDSN(username, password, host string, options Options) string
	GetCredentials(ctx context.Context, clusterName, username string, options Options) (credentials.GetPostgresCredentialsResponse, error)
	DownloadSSLRootCert(ctx context.Context, clusterName string, options Options) error
	GetEndpoint(clusterName string, options Options) string
}

type Options struct {
	Database        string
	Port            string
	Namespace       string
	Host            string // If host is not specified it will work with clusterName
	Role            string // master or replica
	SSLRootCertPath string
	SSLMode         string
	SSLDownload     bool
}

type PG struct {
	CredentialsProvider credentials.Credentials
}

// AutoSetup This is optional to run if you want to setup automatically and avoid piping everything
func (pg *PG) AutoSetup(env string) error {
	env, err := environment.DetermineEnvironment(env)
	if err != nil {
		return err
	}

	factory := credentials.Factory{Providers: map[string]credentials.Credentials{
		environment.Kubernetes: &credentials.Kubernetes{},
		environment.VM:         &credentials.VM{},
	}}

	pg.CredentialsProvider = factory.Get(env)
	return nil
}

func (pg *PG) GetConnection(ctx context.Context, clusterName, username string, options Options) (*sqlx.DB, error) {
	if options.Host == "" {
		options.Host = pg.GetEndpoint(clusterName, options)
	}

	resp, err := pg.GetCredentials(ctx, clusterName, username, options)
	if err != nil {
		return nil, fmt.Errorf("could not GetCredentials: %v", err)
	}

	if options.SSLDownload {
		if err := pg.DownloadSSLRootCert(ctx, clusterName, options); err != nil {
			return nil, fmt.Errorf("could not DownloadSSLRootCert: %v", err)
		}
	}

	dsn := pg.GetPostgresDSN(resp.Username, resp.Password, options.Host, options)

	return pg.Connect(dsn)
}

func (pg *PG) Connect(dsn string) (*sqlx.DB, error) {
	conn, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	conn.SetMaxIdleConns(1)
	conn.SetMaxOpenConns(1)
	conn.SetConnMaxLifetime(0)
	return conn, nil
}

func (pg *PG) GetPostgresDSN(username, password, host string, options Options) string {
	if options.Port == "" {
		options.Port = "5432"
	}
	if options.Database == "" {
		options.Database = "postgres"
	}
	if options.SSLRootCertPath == "" && options.SSLMode == "" {
		options.SSLMode = "disable"
	} else if options.SSLRootCertPath != "" && options.SSLMode == "" {
		options.SSLMode = "verify-ca"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(
		"postgresql://%v:%v@%v:%v/%v?sslmode=%v",
		username,
		password,
		host,
		options.Port,
		options.Database,
		options.SSLMode,
	))
	if options.SSLRootCertPath != "" {
		sb.WriteString(fmt.Sprintf("&sslrootcert=%v", options.SSLRootCertPath))
	}

	return sb.String()
}

func (pg *PG) GetCredentials(ctx context.Context, clusterName, username string, options Options) (credentials.GetPostgresCredentialsResponse, error) {
	postgresCredentials, err := pg.CredentialsProvider.GetPostgresCredentials(
		ctx,
		clusterName,
		username,
		credentials.Options{Namespace: options.Namespace},
	)
	if err != nil {
		return credentials.GetPostgresCredentialsResponse{}, fmt.Errorf("could not GetPostgresCredentials: %v", err)
	}

	return postgresCredentials, nil
}

func (pg *PG) DownloadSSLRootCert(ctx context.Context, clusterName string, options Options) error {
	cert, err := pg.CredentialsProvider.GetPostgresSSLRootCert(ctx, clusterName, credentials.Options{Namespace: options.Namespace})
	if err != nil {
		return fmt.Errorf("could not GetPostgresSSLRootCert: %v", err)
	}

	return os.WriteFile(options.SSLRootCertPath, cert.RootCertBytes, 0777)
}

func (pg *PG) GetEndpoint(clusterName string, options Options) string {
	if options.Role == "replica" {
		return fmt.Sprintf("%s-repl.%s.svc.%s", clusterName, options.Namespace, "cluster.local")
	}

	return fmt.Sprintf("%s.%s.svc.%s", clusterName, options.Namespace, "cluster.local")
}
