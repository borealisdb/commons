package postgresql

import (
	"context"
	"github.com/borealisdb/commons/credentials"
	"github.com/borealisdb/commons/mocks"
	"github.com/golang/mock/gomock"
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
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "default dsn",
			fields: fields{CredentialsProvider: nil, ClusterName: "mycluster"},
			args: args{
				password: "123",
				options:  Options{},
			},
			want: "postgresql://postgres:123@mycluster.default.svc.cluster.local:5432/postgres?sslmode=disable",
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
					Username:        "admin",
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
			ctrl := gomock.NewController(t)
			mockCredentials := mocks.NewMockCredentials(ctrl)
			if tt.args.options.Host == "" {
				mockCredentials.EXPECT().
					GetClusterEndpoint(context.Background(), tt.fields.ClusterName, tt.args.options.Role).
					Return(credentials.GetClusterEndpointResponse{Endpoint: "mycluster.default.svc.cluster.local"}, nil)
			}
			tt.fields.CredentialsProvider = mockCredentials
			pg := PG{CredentialsProvider: mockCredentials}
			got, err := pg.GetDSN(tt.fields.ClusterName, tt.args.password, tt.args.options)
			if err != nil {
				t.Errorf("GetCredentials() error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("GetDSN() = %v, want %v", got, tt.want)
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
				options: Options{
					Username: "admin",
				},
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
						Endpoint: "mycluster.default.svc.cluster.local",
					}, nil)
			}
			mockCredentials.EXPECT().
				GetPostgresCredentials(context.Background(), tt.fields.clusterName, tt.args.username, tt.fields.credOptions).
				Return(tt.want, nil)

			pg := PG{
				CredentialsProvider: mockCredentials,
			}
			got, err := pg.GetCredentials(context.Background(), tt.fields.clusterName, tt.fields.options)
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
