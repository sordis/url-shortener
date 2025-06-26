package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMustLoadConfig(t *testing.T) {
	t.Run("successful config load", func(t *testing.T) {
		configContent := `
env: test
storage_path: "/tmp/test.db"
http_server:
  address: ":8080"
  timeout: 5s
  idle_timeout: 30s
  user: "admin"
  password: "testpass"
`
		tmpFile, err := os.CreateTemp("", "test_config_*.yml")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.WriteString(configContent)
		require.NoError(t, err)
		tmpFile.Close()

		t.Setenv("CONFIG_PATH", tmpFile.Name())
		t.Setenv("HTTP_SERVER_PASSWORD", "envpass") // Проверяем приоритет переменных окружения

		cfg := MustLoadConfig()

		// Проверяем значения
		require.Equal(t, "test", cfg.Env)
		require.Equal(t, "/tmp/test.db", cfg.StoragePath)
		require.Equal(t, ":8080", cfg.HTTPServer.Address)
		require.Equal(t, 5*time.Second, cfg.HTTPServer.Timeout)
		require.Equal(t, 30*time.Second, cfg.HTTPServer.IdleTimeout)
		require.Equal(t, "admin", cfg.HTTPServer.User)
		require.Equal(t, "envpass", cfg.HTTPServer.Password)
	})

	t.Run("missing config file", func(t *testing.T) {
		t.Setenv("CONFIG_PATH", "nonexistent.yml")

		
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("MustLoadConfig should panic on missing file")
			}
		}()

		_ = MustLoadConfig()
	})

	t.Run("missing required env", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "test_config_*.yml")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		// Конфиг без обязательного поля storage_path
		_, err = tmpFile.WriteString(`
env: test
http_server:
  user: "admin"
  password: "testpass"
`)
		require.NoError(t, err)
		tmpFile.Close()

		t.Setenv("CONFIG_PATH", tmpFile.Name())

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("MustLoadConfig should panic on invalid config")
			}
		}()

		_ = MustLoadConfig()
	})
}

func TestConfig_EnvironmentVariables(t *testing.T) {
	t.Run("override with env variables", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "test_config_*.yml")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.WriteString(`
env: "should_be_overridden"
storage_path: "/tmp/test.db"
http_server:
  address: ":8080"
  user: "admin"
  password: "testpass"
`)
		require.NoError(t, err)
		tmpFile.Close()

		t.Setenv("CONFIG_PATH", tmpFile.Name())
		t.Setenv("ENV", "prod")
		t.Setenv("HTTP_TIMEOUT", "10s")

		cfg := MustLoadConfig()

		require.Equal(t, "prod", cfg.Env)
		require.Equal(t, 10*time.Second, cfg.HTTPServer.Timeout) 
	})
}
