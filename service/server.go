package service

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

// NewServer configures and returns a Server.
func NewServer() *negroni.Negroni {

	formatter := render.New(render.Options{
		IndentJSON: true,
	})

	n := negroni.Classic()
	mx := mux.NewRouter()

	initRoutes(mx, formatter)

	n.UseHandler(mx)
	return n
}

func initRoutes(mx *mux.Router, formatter *render.Render) {
	mx.HandleFunc("/{act}/{id}/{time}", testHandler(formatter)).Methods("GET")
}

func testHandler(formatter *render.Render) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id := vars["id"]
		act := vars["act"]
		time := vars["time"]
		k, _ := strconv.Atoi(time)
		formatter.JSON(w, http.StatusOK, struct{ Repeate string }{time})
		for i := 0; i < k; i++ {
			formatter.JSON(w, http.StatusOK, struct{ Test string }{act + " " + id})
		}
	}
}
