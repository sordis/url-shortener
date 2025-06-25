package delete_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"url-shortener/internal/http-server/handlers/url/delete"
	deleteMocks "url-shortener/internal/http-server/handlers/url/delete/mocks"
	"url-shortener/internal/lib/logger/handlers/slogdiscard"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteHandler(t *testing.T) {
	cases := []struct {
		name         string
		alias        string
		mockError    error
		expectedCode int
		expectedBody string
	}{

		{
			name:         "Success",
			alias:        "valid_alias",
			mockError:    nil,
			expectedCode: http.StatusNoContent,
			expectedBody: "",
		},
		{
			name:         "URL not found",
			alias:        "not_found_alias",
			mockError:    storage.ErrURLNotFound,
			expectedCode: http.StatusOK,
			expectedBody: `{"error":"internal server error", "status":"Error"}`,
		},
		{
			name:         "Storage error",
			alias:        "error_alias",
			mockError:    errors.New("internal server error"),
			expectedCode: http.StatusOK,
			expectedBody: `{"error":"failed to get url", "status":"Error"}`,
		},
		{
			name:         "Alias with special chars",
			alias:        "test@alias$123",
			mockError:    nil,
			expectedCode: http.StatusNoContent,
			expectedBody: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			urlRemoverMock := deleteMocks.NewMockURLRemover(t)

			// Настраиваем мок только если alias не пустой
			if tc.alias != "" {
				urlRemoverMock.On("DeleteURL", tc.alias).
					Return(tc.mockError).
					Once()
			}

			handler := delete.New(slogdiscard.NewDiscardLogger(), urlRemoverMock)

			// Создаем роутер и тестовый сервер
			r := chi.NewRouter()
			r.Delete("/{alias}", handler)

			req, err := http.NewRequest(http.MethodDelete, "/"+tc.alias, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedCode, rr.Code)

			if tc.expectedBody != "" {
				assert.JSONEq(t, tc.expectedBody, rr.Body.String())
			} else {
				assert.Empty(t, rr.Body.String())
			}
		})
	}
}
