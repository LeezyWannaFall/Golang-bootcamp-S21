package main

import (
    "TicTacToe/internal/di"
    "TicTacToe/internal/web"
    "context"
    "net/http"
    "github.com/go-chi/chi/v5"
    "go.uber.org/fx"
)

func main() {
    fx.New(
        di.AppModule,
        fx.Invoke(StartServer),
    ).Run()
}

func StartServer(lc fx.Lifecycle, h *web.Handler) {
	r := chi.NewRouter()
	r.Post("/game", h.StartGame)
    r.Post("/game/{id}", h.NextMove)

    server := &http.Server{
        Addr:    ":8080",
        Handler: r,
    }

    lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
            // запускаем в горутине, чтобы не блокировать старт fx
            go func() {
                server.ListenAndServe();
            }()
            return nil
        },
        OnStop: func(ctx context.Context) error {
            return server.Shutdown(ctx)
        },
    })
}