package html

import (
	"net/http"

	"github.com/gorilla/pat"
	"github.com/jansemmelink/log"
)

//Router is created for each module, e.g. for all the /auth/... links
//with:	NewRouter("auth").
//			Post("register", register).
//			Post("login", login)
//Where the functions register() and login() implements the
//			HandlerFunc

//HandlerFunc ...
type HandlerFunc func(res ISessionRes, req *http.Request) (string, error)

//NewRouter ...
//	name must be simple, e.g. "auth"
func NewRouter(name string) *Router {
	return &Router{
		Router: pat.New(),
		name:   name,
	}
}

//Router ...
type Router struct {
	*pat.Router
	name string
}

//Form get a template and post form values then calls a user handler ...
func (sr *Router) Form(name string, postHandler HandlerFunc) {
	path := "/" + sr.name + "/" + name
	tmpl := "./" + sr.name + "/templates/" + name + ".html"
	log.Debugf("Adding Form(%s) path=%s tmpl=%s", name, path, tmpl)
	//todo: parse template now and fail on error!

	//define http.GET route to display the form
	sr.Router.Get(
		path,
		func(res http.ResponseWriter, req *http.Request) {
			log.Debugf("Displaying Form Template")
			//get the session or create a new one
			session := NewSession()
			//call user handler with a sessionRes() instead of normal http.ResponseWriter
			sessionRes{res, session}.FormTemplate(tmpl, copyURLParams(req.URL.Query()))
		})

	//define http.POST route to process posted form values
	sr.Router.Post(
		path,
		func(res http.ResponseWriter, req *http.Request) {
			log.Debugf("Calling POST Handler...")
			//get the session or create a new one
			session := NewSession()
			//call user handler with a sessionRes() instead of normal http.ResponseWriter
			newURL, err := postHandler(sessionRes{res, session}, req)
			if err != nil {
				//show error
				log.Errorf("ERROR in postHandler: %v", err)
				sessionRes{res, session}.Message("Error", "%s", err)
			} else {
				//redirect
				if len(newURL) > 0 {
					log.Debugf("ROUTER to %s ...", newURL)
					http.Redirect(res, req, newURL, http.StatusOK)
				} else {
					log.Debugf("NOT REDIRECT")
				}
			}
		})
}

//Post ...
func (sr *Router) Post(name string, hf HandlerFunc) {
	sr.Router.Post(
		"/"+sr.name+"/"+name,
		func(res http.ResponseWriter, req *http.Request) {
			log.Debugf("Calling Handler...")
			//get the session or create a new one
			session := NewSession()
			//call user handler with a sessionRes() instead of normal http.ResponseWriter
			hf(sessionRes{res, session}, req)
		})
}
