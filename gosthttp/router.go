package gosthttp

import (
	"fmt"
	"net/http"

	"github.com/geodan/gost/sensorthings"
	"github.com/gorilla/mux"
)

// NewRouter creates a new mux.Router and sets up all endpoints defind in the sensothings api
func NewRouter(api *sensorthings.SensorThingsAPI) *mux.Router {
	// Note: tried julienschmidt/httprouter instead of gorilla/mux but had some
	// problems with interfering endpoints cause of the wildcard used for the (id) in requests
	a := *api
	endpoints := *a.GetEndpoints()
	router := mux.NewRouter().StrictSlash(true)
	router.Path("/").Handler(http.FileServer(http.Dir("./gostsite/")))

	for _, endpoint := range endpoints {
		ep := endpoint
		for _, op := range ep.Operations {
			operation := op
			method := fmt.Sprintf("%s", operation.OperationType)
			if operation.Handler == nil {
				continue
			}

			router.Methods(method).
				Path(operation.Path).
				HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					operation.Handler(w, r, &ep, api)
				})
		}
	}

	return router
}
