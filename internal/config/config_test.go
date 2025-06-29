package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMustLoadConfig(t *testing.T) {
	// Сохраняем оригинальные значения переменных окружения
	originalConfigPath := os.Getenv("CONFIG_PATH")
	originalPassword := os.Getenv("HTTP_SERVER_PASSWORD")

	t.Cleanup(func() {
		// Восстанавливаем оригинальные значения после теста
		err := os.Setenv("CONFIG_PATH", originalConfigPath)
		require.NoError(t, err)
		errP := os.Setenv("HTTP_SERVER_PASSWORD", originalPassword)
		require.NoError(t, errP)
	})

	t.Run("successful config load", func(t *testing.T) {
		// Создаём временный файл
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yml")

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
		err := os.WriteFile(configPath, []byte(configContent), 0644)
		require.NoError(t, err)

		errConf := os.Setenv("CONFIG_PATH", configPath)
		require.NoError(t, errConf)
		errPass := os.Setenv("HTTP_SERVER_PASSWORD", "envpass")
		require.NoError(t, errPass)

		cfg := MustLoadConfig()

		require.Equal(t, "test", cfg.Env)
		require.Equal(t, "/tmp/test.db", cfg.StoragePath)
		require.Equal(t, ":8080", cfg.HTTPServer.Address)
		require.Equal(t, "envpass", cfg.HTTPServer.Password)
	})
}