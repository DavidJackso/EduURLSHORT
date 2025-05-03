package delete

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	resp "urlshorter/internal/lib/api/response"
)

type URLDelete interface {
	Delete(alias string) error
}

func New(log *slog.Logger, urlDelete URLDelete) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Error("alias is empty")
			render.JSON(w, r, resp.Error("alias is required"))
			return
		}

		log.Info("attempting to delete url", slog.String("alias", alias))

		err := urlDelete.Delete(alias)
		if err != nil {
			log.Error("failed to delete url", slog.String("alias", alias), slog.Any("error", err))
			render.JSON(w, r, resp.Error("failed to delete url"))
			return
		}

		log.Info("url successfully deleted", slog.String("alias", alias))
		render.JSON(w, r, resp.OK())
	}
}
