package user

import (
	"net/http"

	html "github.com/jansemmelink/trotsek/lib/html"
)

//Router to include auth functions in web site
func Router() http.Handler {
	r := html.NewRouter("user")
	//	r.Form("home", )
	return r
}

func userHome(res html.ISessionRes, req *http.Request) {
	res.Message("User's Home",
		"Soon to come...")
}
