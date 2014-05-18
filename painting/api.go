package painting

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"fmt"
	"github.com/danielphan/dandubois/apiutil"
	"github.com/gorilla/mux"
	"net/http"
)

var router *mux.Router

func Routes(r *mux.Router) {
	r.Path("/categories").HandlerFunc(apiutil.Wrap(listCategories))
	r.Path("/media").HandlerFunc(apiutil.Wrap(listMedia))

	ps := r.Path("/paintings").Subrouter()
	ps.Methods("GET").HandlerFunc(apiutil.Wrap(listPaintings))
	ps.Methods("POST").HandlerFunc(apiutil.Wrap(createPainting))

	p := r.Path("/paintings/{ID:[0-9]+}").Subrouter()
	p.Methods("GET").HandlerFunc(apiutil.Wrap(showPainting))
	p.Methods("PUT").HandlerFunc(apiutil.Wrap(editPainting))

	router = r
}

func listCategories(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	q := datastore.NewQuery(Kind).
		Project("Categories").
		Distinct().
		Order("Categories")

	var paintings []*struct{ Categories []string }
	_, err := q.GetAll(c, &paintings)
	if err != nil {
		return err
	}

	var categories []string
	for _, p := range paintings {
		categories = append(categories, p.Categories[0])
	}

	err = json.NewEncoder(w).Encode(categories)
	if err != nil {
		return err
	}

	return nil
}

func listMedia(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	q := datastore.NewQuery(Kind).
		Project("Media").
		Distinct().
		Order("Media")

	var paintings []*struct{ Media []string }
	_, err := q.GetAll(c, &paintings)
	if err != nil {
		return err
	}

	var media []string
	for _, p := range paintings {
		media = append(media, p.Media[0])
	}

	err = json.NewEncoder(w).Encode(media)
	if err != nil {
		return err
	}

	return nil
}

func listPaintings(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	q := datastore.NewQuery(Kind).
		Order("-Year").
		Order("-ID")

	var paintings []*Painting
	_, err := q.GetAll(c, &paintings)
	if err != nil {
		return err
	}

	err = json.NewEncoder(w).Encode(paintings)
	if err != nil {
		return err
	}

	return nil
}

func createPainting(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)

	var p Painting
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		return err
	}

	err = p.Save(c)
	if err != nil {
		return err
	}

	return nil
}

func showPainting(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	id := mux.Vars(r)["ID"]
	p, err := Get(c, ID(id))
	if err != nil {
		return err
	}

	err = json.NewEncoder(w).Encode(p)
	if err != nil {
		return err
	}

	return nil
}

func editPainting(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprintln(w, "editPainting")
	return nil
}
