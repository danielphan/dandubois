package migrate

import (
	"appengine"
	"appengine/blobstore"
	"appengine/datastore"
	"appengine/urlfetch"
	"github.com/danielphan/ae/apiutil"
	"github.com/danielphan/dandubois/painting"
	"github.com/gorilla/mux"
	"net/http"
)

const baseURL = "http://static.danieldubois.net/"

func Routes(r *mux.Router) {
	r = r.PathPrefix("/migrate").Subrouter()
	r.Methods("POST").Path("/load").Handler(load)
	r.Methods("POST").Path("/clean").Handler(clean)
}

var load = apiutil.Error(apiutil.Json(apiutil.Admin(
	func (w http.ResponseWriter, r *http.Request) error {
		c := appengine.NewContext(r)
		client := urlfetch.Client(c)

		cs, err := fetchCategories(client)
		if err != nil {
			return err
		}

		ms, err := fetchMedia(client)
		if err != nil {
			return err
		}

		ps, err := fetchPaintings(client)
		if err != nil {
			return err
		}

		ps.convert(w, c, cs, ms)

		return nil
	})))

var clean = apiutil.Error(apiutil.Json(apiutil.Admin(
	func (w http.ResponseWriter, r *http.Request) error {
		c := appengine.NewContext(r)
		q := datastore.NewQuery(painting.Kind).Project("Image")

		var paintings []*painting.Painting
		keys, err := q.GetAll(c, &paintings)
		if err != nil {
			return err
		}

		var blobKeys []appengine.BlobKey
		for _, p := range paintings {
			if p.Image.BlobKey != appengine.BlobKey(" ") {
				blobKeys = append(blobKeys, p.Image.BlobKey)
			}
		}
		err = blobstore.DeleteMulti(c, blobKeys)
		if err != nil {
			return err
		}

		datastore.DeleteMulti(c, keys)
		if err != nil {
			return err
		}

		return nil
	})))
