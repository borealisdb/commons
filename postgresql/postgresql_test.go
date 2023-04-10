package postgresql

import (
	"github.com/borealisdb/commons/credentials"
	"testing"
)

func TestPG_GetPostgresDSN(t *testing.T) {
	type fields struct {
		CredentialsProvider credentials.Credentials
	}
	type args struct {
		username string
		password string
		host     string
		options  Options
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{name: "default dsn", fields: fields{CredentialsProvider: nil}, args: args{
			username: "admin",
			password: "123",
			host:     "localhost",
			options:  Options{},
		}, want: "postgresql://admin:123@localhost:5432/postgres?sslmode=disable"},
		{name: "dsn with ssl enabled", fields: fields{CredentialsProvider: nil}, args: args{
			username: "admin",
			password: "123",
			host:     "localhost",
			options: Options{
				SSLRootCertPath: "/borealis/root.crt",
			},
		}, want: "postgresql://admin:123@localhost:5432/postgres?sslmode=verify-ca&sslrootcert=/borealis/root.crt"},
		{name: "dsn with custom options", fields: fields{CredentialsProvider: nil}, args: args{
			username: "admin",
			password: "123",
			host:     "localhost",
			options: Options{
				Database:        "users",
				Port:            "5001",
				SSLRootCertPath: "/borealis/root.crt",
				SSLMode:         "verify-full",
			},
		}, want: "postgresql://admin:123@localhost:5001/users?sslmode=verify-full&sslrootcert=/borealis/root.crt"},
		{name: "dsn with custom options", fields: fields{CredentialsProvider: nil}, args: args{
			username: "admin",
			password: "123",
			host:     "localhost",
			options: Options{
				Database: "users",
				Port:     "5001",
				SSLMode:  "require",
			},
		}, want: "postgresql://admin:123@localhost:5001/users?sslmode=require"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pg := &PG{
				CredentialsProvider: tt.fields.CredentialsProvider,
			}
			if got := pg.GetPostgresDSN(tt.args.username, tt.args.password, tt.args.host, tt.args.options); got != tt.want {
				t.Errorf("GetPostgresDSN() = %v, want %v", got, tt.want)
			}
		})
	}
}
