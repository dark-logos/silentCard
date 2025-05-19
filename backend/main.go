package main

   import (
   	"encoding/json"
   	"fmt"
   	"html"
   	"log"
   	"net/http"
   	"os"
   	"time"
   )

   func main() {
   	http.HandleFunc("/save-wish", saveWishHandler)
   	port := os.Getenv("PORT")
   	if port == "" {
   		port = "10000" // Render default
   	}
   	log.Printf("Server starting on port %s", port)
   	log.Fatal(http.ListenAndServe(":"+port, nil))
   }

   type WishRequest struct {
   	Wish string `json:"wish"`
   }

   type WishResponse struct {
   	Success bool   `json:"success"`
   	Error   string `json:"error,omitempty"`
   }

   func saveWishHandler(w http.ResponseWriter, r *http.Request) {
   	// CORS
   	w.Header().Set("Access-Control-Allow-Origin", "https://your-username.github.io")
   	w.Header().Set("Access-Control-Allow-Methods", "POST")
   	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

   	if r.Method == "OPTIONS" {
   		w.WriteHeader(http.StatusOK)
   		return
   	}

   	if r.Method != http.MethodPost {
   		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
   		return
   	}

   	var req WishRequest
   	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
   		sendJSONResponse(w, WishResponse{Success: false, Error: "Invalid JSON"})
   		return
   	}

   	if req.Wish == "" {
   		sendJSONResponse(w, WishResponse{Success: false, Error: "Желание не указано"})
   		return
   	}

   	// Санитизация
   	sanitizedWish := html.EscapeString(req.Wish)

   	// Запись в файл
   	timestamp := time.Now().Format("2006-01-02 15:04:05")
   	entry := fmt.Sprintf("[%s] %s\n", timestamp, sanitizedWish)

   	file, err := os.OpenFile("wishes.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
   	if err != nil {
   		log.Printf("Error opening file: %v", err)
   		sendJSONResponse(w, WishResponse{Success: false, Error: "Не удалось сохранить желание"})
   		return
   	}
   	defer file.Close()

   	if _, err := file.WriteString(entry); err != nil {
   		log.Printf("Error writing to file: %v", err)
   		sendJSONResponse(w, WishResponse{Success: false, Error: "Не удалось сохранить желание"})
   		return
   	}

   	sendJSONResponse(w, WishResponse{Success: true})
   }

   func sendJSONResponse(w http.ResponseWriter, resp WishResponse) {
   	w.Header().Set("Content-Type", "application/json")
   	if err := json.NewEncoder(w).Encode(resp); err != nil {
   		log.Printf("Error encoding JSON: %v", err)
   	}
   }