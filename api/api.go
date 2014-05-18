package api

import (
	"github.com/danielphan/dandubois/painting"
	"github.com/danielphan/dandubois/painting/migrate"
	"github.com/gorilla/mux"
	"net/http"
)

func init() {
	r := mux.NewRouter().
		PathPrefix("/api").
		Subrouter()

	painting.Routes(r)
	migrate.Routes(r)

	http.Handle("/", r)
}
