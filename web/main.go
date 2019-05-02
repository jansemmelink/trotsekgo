package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/pat"
	"html/template"
	"github.com/jansemmelink/log"
)

func main() {
	log.DebugOn()
	if err := http.ListenAndServe("localhost:8060", app()); err != nil {
		panic(fmt.Sprintf("Failed in http server: %v", err))
	}
	fmt.Printf("Terminated normally\n")
}

func app() http.Handler {
	r := pat.New()
	r.Add(http.MethodGet, "/resources", http.StripPrefix("/resources/", http.FileServer(http.Dir("./resources"))))
	r.Get("/", page)
	return r
}

//page from template
func page(res http.ResponseWriter, req *http.Request) {
	data := map[string]string{}
	name:=req.URL.Path
	if name == "/" {
		name = "/home"
	}
	t, err := template.ParseFiles(
		"./templates"+name+".html",
		"./templates/pageHeader.html",
		"./templates/pageFooter.html",
	)
	if err != nil {
		//not found
		t, err = template.ParseFiles(
			"./templates/notFound.html",
			"./templates/pageHeader.html",
			"./templates/pageFooter.html",
		)
		if err != nil {
			panic(log.Wrapf(err, "failed to parse template"))
		}
	}

	err = t.Execute(res, data)
	if err != nil {
		panic(log.Wrapf(err, "failed to execute template"))
	}
}
