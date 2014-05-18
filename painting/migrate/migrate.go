package migrate

import (
	"appengine"
	"appengine/blobstore"
	"appengine/datastore"
	"appengine/urlfetch"
	"github.com/danielphan/dandubois/painting"
	"github.com/gorilla/mux"
	"net/http"
)

const baseURL = "http://static.danieldubois.net/"

func Routes(r *mux.Router) {
	r.Path("/migrate/load").
		HandlerFunc(load)

	r.Path("/migrate/clean").
		HandlerFunc(clean)
}

func load(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	client := urlfetch.Client(c)

	cs, err := fetchCategories(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ms, err := fetchMedia(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ps, err := fetchPaintings(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ps.convert(w, c, cs, ms)
}

func clean(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	q := datastore.NewQuery(painting.Kind).Project("Image")

	var paintings []*painting.Painting
	keys, err := q.GetAll(c, &paintings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var blobKeys []appengine.BlobKey
	for _, p := range paintings {
		if p.Image.BlobKey != appengine.BlobKey(" ") {
			blobKeys = append(blobKeys, p.Image.BlobKey)
		}
	}
	err = blobstore.DeleteMulti(c, blobKeys)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	datastore.DeleteMulti(c, keys)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
