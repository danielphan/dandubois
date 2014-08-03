package painting

import (
	"appengine"
	"appengine/user"
	"encoding/json"
	"errors"
	"github.com/danielphan/ae/apiutil"
	"github.com/danielphan/ae/logger"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var ErrMustLogIn = errors.New("Must be logged in.")

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

	r.Path("/paintings/{ID:[0-9]+}/rotate").
		Subrouter().
		Methods("POST").
		HandlerFunc(apiutil.Wrap(rotatePainting))

	router = r
}

func listCategories(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	categories, err := GetAllCategories(c)
	if err != nil {
		return err
	}

	err = json.NewEncoder(w).Encode(categories)
	if err != nil {
		return err
	}

	return nil
}

func listMedia(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	media, err := GetAllMedia(c)
	if err != nil {
		return err
	}

	err = json.NewEncoder(w).Encode(media)
	if err != nil {
		return err
	}

	return nil
}

func listPaintings(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	paintings, err := GetAll(c)
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
	if !user.IsAdmin(c) {
		s := logger.Error(c, ErrMustLogIn)
		http.Error(w, s, http.StatusUnauthorized)
		return nil
	}

	p := &Painting{}
	err := json.NewDecoder(r.Body).Decode(p)
	if err != nil {
		return err
	}

	err = p.Save(c)
	if err != nil {
		return err
	}

	err = json.NewEncoder(w).Encode(p)
	if err != nil {
		return err
	}

	return nil
}

func get(r *http.Request) (*Painting, error) {
	c := appengine.NewContext(r)
	id := mux.Vars(r)["ID"]
	p, err := Get(c, ID(id))
	if err != nil {
		return nil, err
	}
	return p, nil
}

func showPainting(w http.ResponseWriter, r *http.Request) error {
	p, err := get(r)
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
	c := appengine.NewContext(r)
	if !user.IsAdmin(c) {
		s := logger.Error(c, ErrMustLogIn)
		http.Error(w, s, http.StatusUnauthorized)
		return nil
	}
	return nil
}

func rotatePainting(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	if !user.IsAdmin(c) {
		s := logger.Error(c, ErrMustLogIn)
		http.Error(w, s, http.StatusUnauthorized)
		return nil
	}

	p, err := get(r)
	if err != nil {
		return err
	}

	angle, err := strconv.ParseFloat(r.FormValue("angle"), 64)
	if err != nil {
		return err
	}

	err = p.rotate(c, angle)
	if err != nil {
		return err
	}

	err = json.NewEncoder(w).Encode(p)
	if err != nil {
		return nil
	}

	return nil
}
