package migrate

import (
	"appengine"
	"appengine/urlfetch"
	"net/http"
)

const baseUrl = "http://static.danieldubois.net/"

func init() {
	http.HandleFunc("/api/migrate", migrate)
}

func migrate(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	client := urlfetch.Client(c)

	cs, err := fetchCategories(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	ms, err := fetchMedia(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	ps, err := fetchPaintings(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	ps.convert(c, cs, ms)
}
