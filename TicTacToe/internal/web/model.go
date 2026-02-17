package handler

import "github.com/google/uuid"

type RequestDTO struct {
	GameField [3][3]int `json:"field"`
}

type ResponseDTO struct {
	GameID    uuid.UUID `json:"id"`
    GameField [3][3]int `json:"field"`
    Winner    int       `json:"winner"`
}