package html

import (
	"html/template"
	"net/http"
	"net/url"
	"path"

	"github.com/gorilla/pat"
	"github.com/jansemmelink/log"
)

//ISite ...
type ISite interface {
	//configure the site:
	//WithFiles() serves static files in the specified directory on the specified path
	WithFiles(url, dir string) ISite
	WithPages(dir string) ISite
	WithModule(methods []string, url string, hdlr http.Handler) ISite
	Router() http.Handler
	NewSession() ISession
	Session(id string) ISession
}

type site struct {
	pagesDir string
	router   *pat.Router
	sessions map[string]ISession
}

//NewSite ...
func NewSite() ISite {
	newSite := &site{}
	newSite.router = pat.New()
	return newSite
} //NewSite()

//usage example: WithFiles("/resources", "./resource")
func (site *site) WithFiles(url, dir string) ISite {
	log.Debugf("Files from: %s -> %s", url, dir)
	site.router.Add(http.MethodGet, url, http.StripPrefix(url+"/", http.FileServer(http.Dir(dir))))
	return site
}

func (site *site) WithPages(dir string) ISite {
	site.pagesDir = dir
	return site
}

//WithModule adds a module router to the site
//usage example: mySite.With("/auth", auth.Router())
func (site *site) WithModule(methods []string, url string, hdlr http.Handler) ISite {
	for _, method := range methods {
		site.router.Add(
			method,
			"/"+url,                       //e.g. "/" + "auth" -> "/auth"
			siteHandler{s: site, h: hdlr}) //wrap user handle for this site
	}
	return site
}

func (site site) Router() http.Handler {
	site.router.Get("/", site.page)
	return site.router
}

func (site *site) NewSession() ISession {
	s := NewSession()
	site.sessions[s.ID()] = s
	return s
}

func (site site) Session(id string) ISession {
	if s, ok := site.sessions[id]; ok {
		return s
	}
	return nil
}

func copyURLParams(p url.Values) map[string]string {
	data := make(map[string]string)
	for n,v := range p {
		data[n] = v[0]
	}
	return data
}

//page handler for static pages named on the URL from template
func (site site) page(res http.ResponseWriter, req *http.Request) {
	data := copyURLParams(req.URL.Query())
	name := req.URL.Path
	if name == "/" {
		name = "/home"
	}
	name = path.Dir(name) + "/templates/" + path.Base(name)
	if name[0:2] == "//" {
		name = name[1:]
	}
	t, err := template.ParseFiles(
		"."+name+".html",
		"./templates/pageHeader.html",
		"./templates/pageFooter.html",
	)
	if err != nil {
		//not found
		data["name"] = name
		t, err = template.ParseFiles(
			"./templates/notFound.html",
			"./templates/pageHeader.html",
			"./templates/pageFooter.html",
		)
		if err != nil {
			panic(log.Wrapf(err, "failed to parse template"))
		}
	}

	log.Debugf("Render %s with %+v", t.Name(), data)
	err = t.Execute(res, data)
	if err != nil {
		panic(log.Wrapf(err, "failed to execute template"))
	}
} //site.page()

//siteHandler wraps user module handlers to implement session management
//around all requests
type siteHandler struct {
	s *site
	h http.Handler
}

func (sh siteHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	//get session
	log.Debugf("session tbd")
	session := NewSession()
	sh.h.ServeHTTP(session.Res(res), req)
}
