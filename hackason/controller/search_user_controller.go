package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"uttc_hackason_be/dao"
	"uttc_hackason_be/usecase"
)

func GetUserController(w http.ResponseWriter, r *http.Request, dao *dao.UserDao) {
	name := r.URL.Query().Get("name")
	if name == "" {
		log.Printf("fail: name is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	users, err := usecase.SearchUserUsecase(name, dao)
	bytes, err := json.Marshal(users)
	if err != nil {
		log.Printf("fail: json.Marshal, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.Write(bytes)
}
