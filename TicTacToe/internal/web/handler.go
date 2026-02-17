package web

import (
	"TicTacToe/internal/domain/service"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct{
	service service.DomainInterface
}

func NewHandler(service service.DomainInterface) *Handler {
	return &Handler{service: service}
}

func (h *Handler) StartGame(w http.ResponseWriter, r *http.Request) {
	game, err := h.service.StartGame()
	if err != nil {
		http.Error(w, "could not start game", http.StatusInternalServerError)
		return
	}

	response := ToResponse(game)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) NextMove(w http.ResponseWriter, r *http.Request) {
	var dto RequestDTO

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	id := chi.URLParam(r, "id")

	game := FromRequest(dto)
	game.ID, err = uuid.Parse(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	PlayerMove := h.service.NextMove(game)
	if !PlayerMove {
		http.Error(w, "invalid move", http.StatusBadRequest)
		return
	}

	// AiMove := h.service.NextMove(game)
	// if !AiMove {
	// 	http.Error(w, "invalid move", http.StatusBadRequest)
	// 	return
	// }

	response := ToResponse(game)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}