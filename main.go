package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

var (
	urls    = make(map[string]string)
	mu      sync.RWMutex
	counter = 0
)

func nextShortCode() string {
	mu.Lock()
	defer mu.Unlock()
	counter++
	return fmt.Sprintf("%d", counter)
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var request struct{
		URL string `json:"url"`
	}

	if err:=json.NewDecoder(r.Body).Decode(&request); err!=nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if request.URL==""{
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	shortCode:=nextShortCode()
	mu.Lock()
	urls[shortCode] = request.URL
	mu.Unlock();

	response:= map[string]string{"short_url": fmt.Sprintf("/%s", shortCode)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response);
}

func redirectHandler(w http.ResponseWriter, r *http.Request){
	shortCode:= r.URL.Path[1:];

	mu.RLock()
	originalUrl, exists:=urls[shortCode];
	mu.RUnlock()

	if !exists {
		http.Error(w, "Short Url Not Found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalUrl, http.StatusFound)

}

func main(){
	http.HandleFunc("/shorten", shortenHandler);
	http.HandleFunc("/", redirectHandler);

	fmt.Println("Starting URL Shortner on :8080")
	http.ListenAndServe(":8080", nil);
}
