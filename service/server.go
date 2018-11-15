package service

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

// NewServer configures and returns a Server.
func NewServer() *negroni.Negroni {

	//采用Json模式输出
	formatter := render.New(render.Options{
		IndentJSON: true,
	})

	n := negroni.Classic()

	//使用gorilla/mux库新建路由匹配
	mx := mux.NewRouter()

	//为路由匹配添加处理函数HandlerFunc
	initRoutes(mx, formatter)

	n.UseHandler(mx)
	return n
}

func initRoutes(mx *mux.Router, formatter *render.Render) {
	mx.HandleFunc("/{act}/{id}/{time}", testHandler(formatter)).Methods("GET")
	mx.HandleFunc("/find/{id}", find).Methods("GET")
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

func find(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	vars := mux.Vars(req)
	id := vars["id"]

	//这个写入到w的是输出到客户端的
	fmt.Fprintf(w, "Find request to "+req.Host+req.URL.Path+"\n")
	fmt.Fprintf(w, "Result: Cannot find user "+id+"\n")
}
