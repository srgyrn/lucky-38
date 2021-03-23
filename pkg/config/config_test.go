package config_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/srgyrn/lucky-38/pkg/config"
)

func TestLoad(t *testing.T) {
	appEnvSetter := func(t *testing.T, appenv string) func() {
		oldAppEnv := os.Getenv("APP_ENV")
		if err := os.Setenv("APP_ENV", appenv); err != nil {
			t.Fatalf("os.Setenv() err: %v", err)
		}
		return func() {
			os.Setenv("APP_ENV", oldAppEnv)
		}
	}

	tests := []struct {
		name    string
		appEnv  string
		want    config.Config
		wantErr bool
	}{
		{
			name:   "load dev",
			appEnv: "development",
			want: config.Config{
				Driver: "mysql",
				Source: "mysql://db_admin:admin321@db/lucky",
			},
			wantErr: false,
		},
		{
			name:   "load test",
			appEnv: "test",
			want: config.Config{
				Driver: "postgres",
				Source: "postgresql://db_admin:admin321@db/lucky_test?sslmode=disable",
			},
			wantErr: false,
		},
		{
			name:    "missing file",
			appEnv:  "unknown_env_name",
			want:    config.Config{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			revertAppEnv := appEnvSetter(t, tt.appEnv)
			defer revertAppEnv()

			got, err := config.Load("testdata")

			if !tt.wantErr && err != nil {
				t.Fatalf("Load() unexpected error: %v", err)
			}

			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() = %v, want %v", got, tt.want)
			}
		})
	}
}
