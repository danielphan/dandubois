package painting

import (
	"appengine"
	"encoding/json"
	"github.com/danielphan/ae/apiutil"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func Routes(r *mux.Router) {
	ps := r.Path("/paintings").Subrouter()
	ps.Methods("GET").Handler(listPaintings)
	ps.Methods("POST").Handler(createPainting)

	r.Methods("GET").Path("/paintings/categories").Handler(listCategories)
	r.Methods("GET").Path("/paintings/media").Handler(listMedia)

	p := r.Path("/paintings/{ID:[0-9]+}").Subrouter()
	p.Methods("GET").Handler(showPainting)
	p.Methods("PUT").Handler(editPainting)

	r.Methods("POST").Path("/paintings/{ID:[0-9]+}/rotate").Handler(rotatePainting)
}

var listCategories = apiutil.Error(apiutil.Json(
	func (w http.ResponseWriter, r *http.Request) error {
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
	}))

var listMedia = apiutil.Error(apiutil.Json(
	func (w http.ResponseWriter, r *http.Request) error {
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
	}))

var listPaintings = apiutil.Error(apiutil.Json(
	func (w http.ResponseWriter, r *http.Request) error {
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
	}))

var createPainting = apiutil.Error(apiutil.Json(apiutil.Admin(
	func (w http.ResponseWriter, r *http.Request) error {
		c := appengine.NewContext(r)

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
	})))

func getPainting(r *http.Request) (*Painting, error) {
	c := appengine.NewContext(r)
	id := mux.Vars(r)["ID"]
	p, err := Get(c, ID(id))
	if err != nil {
		return nil, err
	}
	return p, nil
}

var showPainting = apiutil.Error(apiutil.Json(
	func (w http.ResponseWriter, r *http.Request) error {
		p, err := getPainting(r)
		if err != nil {
			return err
		}

		err = json.NewEncoder(w).Encode(p)
		if err != nil {
			return err
		}

		return nil
	}))

var editPainting = apiutil.Error(apiutil.Json(apiutil.Admin(
	func (w http.ResponseWriter, r *http.Request) error {
		return nil
	})))

var rotatePainting = apiutil.Error(apiutil.Json(apiutil.Admin(
	func (w http.ResponseWriter, r *http.Request) error {
		c := appengine.NewContext(r)

		p, err := getPainting(r)
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
	})))
