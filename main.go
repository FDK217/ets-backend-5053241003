package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Server PinjamKampus berjalan pada port 8080...")
	err := http.ListenAndServe(":8080", NewHandler())
	if err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}