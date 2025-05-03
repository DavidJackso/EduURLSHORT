package save

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	resp "urlshorter/internal/lib/api/response"
	"urlshorter/internal/lib/random"
)

// TODO:need move in config
const aliasLength = 6

type URLSaver interface {
	SaveURL(urlToSave string, alias string) error
}

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}
type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to parse request", "error", err)

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decode", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			log.Error("failed to validate request", "error", err)
			//TODO: Refactor this. It returns a bad format error.
			render.JSON(w, r, resp.Error("failed to validate request"))

			return
		}
		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}
		err = urlSaver.SaveURL(req.URL, alias)
		if err != nil {
			log.Info("url already exists", "url", req.URL)

			render.JSON(w, r, resp.Error("url already exists"))

			return
		}
		log.Info("url saved", "url", req.URL)
		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias:    alias,
		})
	}
}
