package tests

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"testing"

	"url-shortener/internal/config"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"

	"url-shortener/internal/http-server/handlers/url/save"
	"url-shortener/internal/lib/api"
	"url-shortener/internal/lib/random"
)

func getTestHost() string {
	cfg := config.MustLoadConfig()
	if cfg.HTTPServer.Address != "" {
		return cfg.HTTPServer.Address
	}
	return "localhost:8080"
}

func TestMain(m *testing.M) {
	// Пробуем разные пути к конфигу
	configPaths := []string{
		"../config/local.yml",
		"config/local.yml",
		"config/prod.yml",
		"../config/prod.yml",
	}

	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			if err := os.Setenv("CONFIG_PATH", path); err != nil {
				log.Fatal("Failed to set CONFIG_PATH:", err)
			}
			break
		}
	}

	code := m.Run()
	os.Exit(code)
}

func TestURLShortener_HappyPath(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   getTestHost(),
	}
	e := httpexpect.Default(t, u.String())
	cfg := config.MustLoadConfig()

	e.POST("/url").
		WithJSON(save.Request{
			URL:   gofakeit.URL(),
			Alias: random.NewRandomString(10),
		}).
		WithBasicAuth(cfg.HTTPServer.User, cfg.HTTPServer.Password).
		Expect().
		Status(200).
		JSON().Object().
		ContainsKey("alias")
}

//nolint:funlen
func TestURLShortener_SaveRedirect(t *testing.T) {
	cfg := config.MustLoadConfig()
	testCases := []struct {
		name  string
		url   string
		alias string
		error string
	}{
		{
			name:  "Valid URL",
			url:   gofakeit.URL(),
			alias: gofakeit.Word() + gofakeit.Word(),
		},
		{
			name:  "Invalid URL",
			url:   "invalid_url",
			alias: gofakeit.Word(),
			error: "invalid URL: URL",
		},
		{
			name:  "Empty Alias",
			url:   gofakeit.URL(),
			alias: "",
		},
		// TODO: add more test cases
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   getTestHost(),
			}

			e := httpexpect.Default(t, u.String())

			// Save

			resp := e.POST("/url").
				WithJSON(save.Request{
					URL:   tc.url,
					Alias: tc.alias,
				}).
				WithBasicAuth(cfg.HTTPServer.User, cfg.HTTPServer.Password).
				Expect().Status(http.StatusOK).
				JSON().Object()

			if tc.error != "" {
				resp.NotContainsKey("alias")

				resp.Value("error").String().IsEqual(tc.error)

				return
			}

			alias := tc.alias

			if tc.alias != "" {
				resp.Value("alias").String().IsEqual(tc.alias)
			} else {
				resp.Value("alias").String().NotEmpty()

				alias = resp.Value("alias").String().Raw()
			}

			// Redirect

			testRedirect(t, alias, tc.url)
		})
	}
}

func testRedirect(t *testing.T, alias string, urlToRedirect string) {
	u := url.URL{
		Scheme: "http",
		Host:   getTestHost(),
		Path:   alias,
	}

	redirectedToURL, err := api.GetRedirect(u.String())
	require.NoError(t, err)

	require.Equal(t, urlToRedirect, redirectedToURL)
}
