package controller

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"uttc_hackason_be/model"
	"uttc_hackason_be/usecase"
)

type TweetController struct {
	TweetUseCase *usecase.TweetUseCase
}

func NewTweetController(db *sql.DB) *TweetController {
	tweetUseCase := usecase.NewTweetUseCase(db)
	return &TweetController{TweetUseCase: tweetUseCase}
}

func (tc *TweetController) MakeTweet(w http.ResponseWriter, r *http.Request) {
	var tweet *model.TweetInfoForHTTPPOST

	if err := json.NewDecoder(r.Body).Decode(&tweet); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := tc.TweetUseCase.MakeTweet(tweet); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	content := map[string]string{
		"uid":     tweet.Uid,
		"content": tweet.Content,
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

func (tc *TweetController) GetTweet(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")

	if page == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	pNum, _ := strconv.Atoi(page)
	tweets, err := tc.TweetUseCase.GetTweet(pNum)
	if err != nil {
		log.Printf("fail: tc.TweetUseCase.GetTweet, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	bytes, err := json.Marshal(tweets)
	if err != nil {
		log.Printf("fail: json.Marshal, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

func (tc *TweetController) RegisterRoute(r *mux.Router) {
	r.HandleFunc("/tweets", tc.MakeTweet).Methods("POST")
	r.HandleFunc("/tweets", tc.GetTweet).Methods("GET")
}
