package redirect

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	resp "urlshorter/internal/lib/api/response"
)

type GetURL interface {
	GetURL(alias string) (string, error)
}

type Request struct {
	Alias string `json:"alias"  validate:"required"`
}

type Response struct {
	URL string `json:"url"`
	resp.Response
}

func New(log *slog.Logger, getURL GetURL) http.HandlerFunc {
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

		url, err := getURL.GetURL(alias)
		if err != nil {
			log.Error("failed to get url", "alias", alias, "error", err)
			render.JSON(w, r, resp.Error("failed to get url"))
			return
		}
		log.Info("get url success", "alias", alias, "url", url)

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)

	}
}
