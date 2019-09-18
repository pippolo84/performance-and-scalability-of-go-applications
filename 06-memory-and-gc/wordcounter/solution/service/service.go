package service

import (
	"encoding/json"
	"net/http"
	"performance-and-scalability-of-go-applications/06-memory-and-gc/wordcounter/solution/social"
	"performance-and-scalability-of-go-applications/06-memory-and-gc/wordcounter/solution/wfreq"
	"strconv"
)

// Search is the http handler for the /search endpoint
func Search(w http.ResponseWriter, r *http.Request) {
	// only GET requests are accepted
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// retrieve user id param from query string
	qparam, ok := r.URL.Query()["user"]
	if !ok || len(qparam) > 1 {
		http.Error(w, "wrong or missing user id", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(qparam[0])
	if err != nil || userID < 1 {
		http.Error(w, "user id must be a positive integer", http.StatusBadRequest)
		return
	}

	// get user name given his ID
	userName, err := social.GetUserName(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// retrieve all posts and comments body related to the user
	userTexts, err := social.GetUserText(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// generate its word frequencies table
	wcounters := wfreq.WordsFreq(userTexts)

	// encode answer as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&wfreq.UserWords{userID, userName, wcounters}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
