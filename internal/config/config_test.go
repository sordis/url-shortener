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
		os.Setenv("CONFIG_PATH", originalConfigPath)
		os.Setenv("HTTP_SERVER_PASSWORD", originalPassword)
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

		// Используем os.Setenv вместо t.Setenv для лучшей совместимости
		os.Setenv("CONFIG_PATH", configPath)
		os.Setenv("HTTP_SERVER_PASSWORD", "envpass")

		cfg := MustLoadConfig()

		require.Equal(t, "test", cfg.Env)
		require.Equal(t, "/tmp/test.db", cfg.StoragePath)
		require.Equal(t, ":8080", cfg.HTTPServer.Address)
		require.Equal(t, "envpass", cfg.HTTPServer.Password)
	})
}