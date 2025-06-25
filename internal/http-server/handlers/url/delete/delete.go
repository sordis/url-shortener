package delete

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
type URLRemover interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlRemover URLRemover) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.url.delete.new"

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

		err := urlRemover.DeleteURL(alias)

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

		log.Info("delete url", slog.String("alias", alias))

		w.WriteHeader(http.StatusNoContent)

	}

}
