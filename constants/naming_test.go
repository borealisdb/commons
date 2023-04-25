package constants

import "testing"

func TestGetClusterEndpoint(t *testing.T) {
	type args struct {
		clusterName string
		namespace   string
		role        string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "get endpoint for master",
			args: args{
				clusterName: "mycluster",
				namespace:   "",
				role:        RoleMaster,
			},
			want: "mycluster.default.svc.cluster.local",
		},
		{
			name: "get endpoint for replica",
			args: args{
				clusterName: "mycluster",
				namespace:   "",
				role:        RoleReplica,
			},
			want: "mycluster-repl.default.svc.cluster.local",
		},
		{
			name: "get endpoint for other namespaces",
			args: args{
				clusterName: "mycluster",
				namespace:   "test",
				role:        RoleReplica,
			},
			want: "mycluster-repl.test.svc.cluster.local",
		},
		{
			name: "get endpoint when no role is specified",
			args: args{
				clusterName: "mycluster",
				namespace:   "",
				role:        "",
			},
			want: "mycluster.default.svc.cluster.local",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetClusterEndpoint(tt.args.clusterName, tt.args.namespace, tt.args.role); got != tt.want {
				t.Errorf("GetClusterEndpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}
