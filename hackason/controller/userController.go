package controller

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"uttc_hackason_be/model"
	"uttc_hackason_be/usecase"
)

type UserController struct {
	UserUsecase *usecase.UserUseCase
}

func NewUserController(db *sql.DB) *UserController {
	userUsecase := usecase.NewUserUseCase(db)
	return &UserController{
		UserUsecase: userUsecase,
	}
}

func (uc *UserController) Register(w http.ResponseWriter, r *http.Request) {
	var user model.UserInfoForHTTPPOST

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("fail: json.NewDecoder.Decode, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := uc.UserUsecase.RegisterUser(&user); err != nil {
		log.Printf("fail: usecase.RegisterUser, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	content := map[string]string{
		"id":   user.ID,
		"name": user.Name,
	}

	bytes, err := json.Marshal(content)
	if err != nil {
		log.Printf("fail: json.Marshal, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.Write(bytes)
}

func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		log.Printf("fail: name is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	users, err := uc.UserUsecase.LoginUser(id)
	bytes, err := json.Marshal(users)
	if err != nil {
		log.Printf("fail: json.Marshal, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.Write(bytes)

}

func (uc *UserController) RegiterRoutes(r *mux.Router) {

	r.HandleFunc("/user", uc.Register).Methods("POST")
	r.HandleFunc("/user", uc.Login).Methods("GET")
}