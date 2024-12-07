package controller

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"uttc_hackason_be/model"
	"uttc_hackason_be/usecase"
)

type LikesController struct {
	LikesUseCase *usecase.LikesUseCase
}

func NewLikesController(db *sql.DB) *LikesController {
	likesUseCase := usecase.NewLikesUseCase(db)
	return &LikesController{LikesUseCase: likesUseCase}
}

func (lc *LikesController) GetLikes(w http.ResponseWriter, r *http.Request) {
	var likeInfo *model.LikeInfoPost

	if err := json.NewDecoder(r.Body).Decode(&likeInfo); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	isLike, countLikes, err := lc.LikesUseCase.ToggleLikes(likeInfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	var response model.LikeResponse = model.LikeResponse{
		IsLike:    isLike,
		CountLike: countLikes,
	}
	bytes, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

func (lc *LikesController) RegisterRoute(r *mux.Router) {
	r.HandleFunc("/likes", lc.GetLikes).Methods("POST")
}
