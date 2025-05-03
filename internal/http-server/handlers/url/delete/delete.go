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
	Delete(alias string) (resp *http.Response)
}

func New(log *slog.Logger, urlDelete URLDelete) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		if alias == "" {
			log.Error("failed to parse alias", "error")

			render.JSON(w, r, resp.Error("failed to parse alias"))

			return
		}

		log.Info("alias parse", slog.Any("alias", alias))

		err := urlDelete.Delete(alias)
		if err != nil {
			log.Error("failed to delete url", "alias", alias, "error", err)
			render.JSON(w, r, resp.Error("failed to delete url"))
			return
		}
		log.Info("url success deleted", "alias", alias)

	}
}
