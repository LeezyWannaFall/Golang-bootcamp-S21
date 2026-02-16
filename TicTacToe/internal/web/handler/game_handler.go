package handler

import (
	"TicTacToe/domain/service"
	"TicTacToe/domain/model"
)

type Handler struct{
	service service.GameService
}

func NewHandler(service service.GameService) *Handler {
	return &Handler{service: service}
}

/*
	NextMove(game *model.Game) (int, int)
	Validate(game *model.Game) error
	IsGameOver(game *model.Game) bool
*/

func (h* Handler) NextMove(game *model.Game) (int, int) {

}

func (h* Handler) Validate(game *model.Game) error {

}

func (h* Handler) IsGameOver(game *model.Game) bool {

}