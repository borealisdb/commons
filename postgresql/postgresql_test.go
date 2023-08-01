package postgresql

import (
	"context"
	"github.com/borealisdb/commons/credentials"
	"github.com/borealisdb/commons/mocks"
	"github.com/golang/mock/gomock"
	"os"
	"reflect"
	"testing"
)

func TestPG_GetPostgresDSN(t *testing.T) {
	type fields struct {
		CredentialsProvider credentials.Credentials
		ClusterName         string
	}
	type args struct {
		password string
		options  Options
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		want     string
		doBefore func()
	}{
		{
			name:   "default dsn",
			fields: fields{CredentialsProvider: credentials.Environment{}, ClusterName: "mycluster"},
			args: args{
				options: Options{},
			},
			want: "postgresql://admin:123@myhost:5432/postgres?sslmode=disable",
			doBefore: func() {
				os.Setenv("mycluster_CLUSTER_HOSTNAME", "myhost")
				os.Setenv("mycluster_CLUSTER_USERNAME", "admin")
				os.Setenv("mycluster_admin_CLUSTER_PASSWORD", "123")
			},
		},

		{
			name:   "dsn with ssl enabled",
			fields: fields{CredentialsProvider: nil, ClusterName: "mycluster"},
			args: args{
				password: "123",
				options: Options{
					SSLRootCertPath: "/borealis/root.crt",
				},
			},
			want: "postgresql://postgres:123@mycluster.default.svc.cluster.local:5432/postgres?sslmode=verify-ca&sslrootcert=/borealis/root.crt",
		},
		{
			name:   "dsn with custom options",
			fields: fields{CredentialsProvider: nil, ClusterName: "mycluster"},
			args: args{
				password: "123",
				options: Options{
					Database:        "users",
					Port:            "5001",
					Host:            "localhost",
					SSLRootCertPath: "/borealis/root.crt",
					SSLMode:         "verify-full",
				},
			},
			want: "postgresql://admin:123@localhost:5001/users?sslmode=verify-full&sslrootcert=/borealis/root.crt",
		},
		{
			name:   "dsn with custom options",
			fields: fields{CredentialsProvider: nil, ClusterName: "mycluster"},
			args: args{
				password: "123",
				options: Options{
					Database: "users",
					Port:     "5001",
					Host:     "localhost",
					SSLMode:  "require",
				},
			},
			want: "postgresql://postgres:123@localhost:5001/users?sslmode=require",
		},
		{
			name:   "dsn with custom Namespace",
			fields: fields{CredentialsProvider: nil, ClusterName: "mycluster"},
			args: args{
				password: "123",
				options: Options{
					SSLMode: "verify-ca",
				},
			},
			want: "postgresql://postgres:123@mycluster.default.svc.cluster.local:5432/postgres?sslmode=verify-ca",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tt.doBefore()

			pg := PG{CredentialsProvider: tt.fields.CredentialsProvider}

			if err := pg.setDefaults(tt.args.options, tt.fields.ClusterName); err != nil {
				t.Errorf("setDefaults() error = %v", err)
				return
			}
			resp, err := pg.getCredentials(ctx, "")
			if err != nil {
				t.Errorf("getCredentials() error = %v", err)
				return
			}

			dsn := pg.getDSN(resp.Password, resp.Username)

			if dsn != tt.want {
				t.Errorf("getDSN() = %v, want %v", dsn, tt.want)
			}
		})
	}
}

func TestPG_GetCredentials(t *testing.T) {
	type fields struct {
		options     Options
		clusterName string
		credOptions credentials.Options
	}
	type args struct {
		username string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    credentials.GetPostgresCredentialsResponse
		wantErr bool
	}{
		{
			name: "default",
			fields: fields{
				options:     Options{},
				clusterName: "mycluster",
				credOptions: credentials.Options{},
			},
			args: args{
				username: "postgres",
			},
			want: credentials.GetPostgresCredentialsResponse{
				Username: "postgres",
				Password: "123",
			},
			wantErr: false,
		},
		{
			name: "different user",
			fields: fields{
				options:     Options{},
				clusterName: "mycluster",
				credOptions: credentials.Options{},
			},
			args: args{
				username: "admin",
			},
			want: credentials.GetPostgresCredentialsResponse{
				Username: "admin",
				Password: "123",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockCredentials := mocks.NewMockCredentials(ctrl)
			if tt.fields.options.Host == "" {
				mockCredentials.EXPECT().
					GetClusterEndpoint(context.Background(), tt.fields.clusterName, tt.fields.options.Role).
					Return(credentials.GetClusterEndpointResponse{
						Hostname: "mycluster.default.svc.cluster.local",
					}, nil)
			}
			mockCredentials.EXPECT().
				GetPostgresCredentials(context.Background(), tt.fields.clusterName, tt.args.username, tt.fields.credOptions).
				Return(tt.want, nil)

			pg := PG{
				CredentialsProvider: mockCredentials,
			}
			got, err := pg.GetCredentials(context.Background(), tt.fields.clusterName, "", tt.fields.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCredentials() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCredentials() got = %v, want %v", got, tt.want)
			}
		})
	}
}
