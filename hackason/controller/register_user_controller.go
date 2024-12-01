package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"uttc_hackason_be/dao"
	"uttc_hackason_be/model"
	"uttc_hackason_be/usecase"
)

func RegisterUserController(w http.ResponseWriter, r *http.Request, dao *dao.UserDao) {
	if err := json.NewDecoder(r.Body).Decode(&model.PostData); err != nil {
		log.Printf("fail: json.NewDecoder.Decode, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	id, err := usecase.RegisterUser(dao)
	if err != nil {
		log.Printf("fail: usecase.RegisterUser, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 結果の出力
	w.WriteHeader(http.StatusOK)
	content := map[string]string{
		"id": id.String(),
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
