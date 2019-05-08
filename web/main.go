package main

import (
	"fmt"
	"net/http"
	"github.com/jansemmelink/log"
	"github.com/jansemmelink/trotsek/web/auth"
	"github.com/jansemmelink/trotsek/web/user"
	"github.com/jansemmelink/trotsek/lib/html"
)

func main() {
	log.DebugOn()
	if err := http.ListenAndServe("localhost:8060", app()); err != nil {
		panic(fmt.Sprintf("Failed in http server: %v", err))
	}
	fmt.Printf("Terminated normally\n")
}

func app() http.Handler {
	site := html.NewSite().
		WithFiles("/resources", "./resources").
		WithPages("./").
		WithModule([]string{http.MethodPost, http.MethodGet}, "auth", auth.Router()).
		WithModule([]string{http.MethodPost, http.MethodGet}, "user", user.Router())
	return site.Router()
}

