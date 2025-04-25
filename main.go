package main

import (
	"fmt"
	"groupieTracker/functions"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", functions.Artists)
	http.HandleFunc("/static/", functions.StaticFiles)
	http.HandleFunc("/api/artists/", functions.GetArtistsData)
	http.HandleFunc("/artist/{id}", functions.ArtistDetails)

	fmt.Println("Server is running on: http://localhost:8085")
	err := http.ListenAndServe(":8085", nil)
	if err != nil {
		log.Fatal(err)
	}

}
