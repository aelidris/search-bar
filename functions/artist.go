package functions

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"text/template"
)

func Artists(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		ErrorPage(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}
	if r.URL.Path != "/" {
		ErrorPage(w, http.StatusNotFound, "Page Not Found")
		return
	}
	tmp, err := template.ParseFiles("./html/artists.html")
	if err != nil {
		ErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	content, err := GetContent(w, artistAPI, "")
	if err != nil {
		ErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	if err := json.Unmarshal(content, &artist); err != nil {
		http.Error(w, "Failed to unmarshal JSON", http.StatusInternalServerError)
		return
	}
	if err := tmp.Execute(w, artist); err != nil {
		ErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
}

// Add this function to your `functions` package
func GetArtistsData(w http.ResponseWriter, r *http.Request) {
	content, err := GetContent(w, artistAPI, "")
	if err != nil {
		ErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	if err := json.Unmarshal(content, &artist); err != nil {
		http.Error(w, "Failed to unmarshal JSON", http.StatusInternalServerError)
		return
	}
	var wg sync.WaitGroup
	wg.Add(len(artist)) // Add the number of artists to the WaitGroup

	// Use a mutex to safely update the artist slice concurrently
	var mu sync.Mutex

	// Fetch location data for each artist concurrently
	for i := 0; i < len(artist); i++ {
		go func(i int) {
			defer wg.Done() // Mark the goroutine as done when it finishes

			// Make sure we're accessing a valid index
			if i >= 0 && i < len(artist) {
				contentLocation, err := GetContent(w, locationAPI, strconv.Itoa(i+1))
				if err != nil {
					ErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
					return
				}
				if contentLocation == nil {
					return // If GetContent fails, it will already call http.Error
				}
				if err := json.Unmarshal(contentLocation, &location); err != nil {
					fmt.Println("Error unmarshaling location data:", err)
					return
				}

				// Safely update the artist slice
				mu.Lock()
				if artist[i].Id == location.Id {
					artist[i].Locations = strings.Join(location.Locations, " ")
				}
				mu.Unlock()
			}
		}(i)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	// Manually write JSON response
	jsonData, err := json.Marshal(artist)
	if err != nil {
		ErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func ArtistDetails(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		ErrorPage(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}
	tmp, err := template.ParseFiles("./html/artistDetails.html")
	if err != nil {
		ErrorPage(w, http.StatusNotFound, "Page Not Found")
		return
	}

	strId := strings.TrimPrefix(r.URL.Path, "/artist/")
	id, err := strconv.Atoi(strId)
	if err != nil {
		ErrorPage(w, http.StatusBadRequest, "Bad Request")
		return
	}
	if id < 1 || id > len(artist) {
		ErrorPage(w, http.StatusBadRequest, "Bad Request")
		return
	}
	contentLocation, err := GetContent(w, locationAPI, strId)
	if err != nil {
		ErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	err = json.Unmarshal(contentLocation, &location)
	if err != nil {
		ErrorPage(w, http.StatusBadRequest, "Bad Request")
		return
	}

	contentDate, err := GetContent(w, dateAPI, strId)
	if err != nil {
		ErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	err = json.Unmarshal(contentDate, &date)
	if err != nil {
		ErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	contentRelation, err := GetContent(w, relationAPI, strId)
	if err != nil {
		ErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	err = json.Unmarshal(contentRelation, &relation)
	if err != nil {
		ErrorPage(w, http.StatusBadRequest, "Bad Request")
		return
	}
	artist[id-1].Locations = location.Locations
	artist[id-1].ConcertDates = strings.Join(date.Dates, " ")
	// artist[id-1].RelationsMap = relation.Relations
	artist[id-1].Relations = relation.Relations

	err = tmp.Execute(w, artist[id-1])
	if err != nil {
		ErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
}

func GetContent(w http.ResponseWriter, API string, strId string) ([]byte, error) {
	response, err := http.Get(API + strId)
	if err != nil {
		http.Error(w, "Failed to fetch artist data", http.StatusInternalServerError)
		return nil, fmt.Errorf("failed to fetch data from API: %v", err)
	}
	defer response.Body.Close()
	content, err := io.ReadAll(response.Body)
	if err != nil {
		http.Error(w, "Failed to read artist data", http.StatusInternalServerError)
		return nil, fmt.Errorf("failed to fetch data from API: %v", err)
	}
	return content, nil
}

func ErrorPage(w http.ResponseWriter, statusCode int, message string) {
	errorData := ErrorData{
		StatusCode: statusCode,
		Message:    message,
	}
	tmpl, err := template.ParseFiles("./html/error.html")
	if err != nil {
		fmt.Println("Error: loading error page")
		return
	}
	// This position of WriteHeader ensures the program only writes to the header once when the error template doesn't exist.
	w.WriteHeader(statusCode)
	err = tmpl.Execute(w, errorData)
	if err != nil {
		http.Error(w, "Error: rendering error page", http.StatusInternalServerError)
		return
	}
}
