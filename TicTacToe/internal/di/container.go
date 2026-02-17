package di

import (
	"TicTacToe/internal/datasource"
	"TicTacToe/internal/domain/service"
	"TicTacToe/internal/web"

	"go.uber.org/fx"
)

var AppModule = fx.Module("app",
    fx.Provide(
        datasource.NewStorage,    // Возвращает *GameStorage
        datasource.NewRepository, // Возвращает DataInterface
        service.NewGameService,   // Возвращает DomainInterface
        web.NewHandler,       // Возвращает *Handler
    ),
)