package functions

import (
	"net/http"
	"os"
)

func StaticFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorPage(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	folders, err := os.ReadDir("static")
	if err != nil {
		ErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	for _, folder := range folders {
		files, _ := os.ReadDir("static/" + folder.Name())
		for _, file := range files {
			if r.URL.Path == "/static/css/"+file.Name() || r.URL.Path == "/static/js/"+file.Name() {
				fs := http.Dir("static")
				http.StripPrefix("/static/", http.FileServer(fs)).ServeHTTP(w, r)
				return
			}
		}
	}

	ErrorPage(w, http.StatusNotFound, "Page Not Found")
}
