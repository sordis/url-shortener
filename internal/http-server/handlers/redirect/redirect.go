package redirect

import (
	"errors"
	"log/slog"
	"net/http"

	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/storage"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

//go:generate mockery
type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.redirect.new"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		resUrl, err := urlGetter.GetURL(alias)

		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("alias not found", "alias", alias)

			render.JSON(w, r, resp.Error("internal server error"))

			return

		}
		if err != nil {
			log.Error("failed to get url")

			render.JSON(w, r, resp.Error("failed to get url"))

			return
		}

		log.Info("got url", slog.String("url", resUrl))

		http.Redirect(w, r, resUrl, http.StatusFound)
	}
}
