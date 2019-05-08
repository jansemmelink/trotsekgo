//Package auth gives auth functions in the web site
package auth

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/jansemmelink/log"
	authlib "github.com/jansemmelink/trotsek/lib/auth"
	html "github.com/jansemmelink/trotsek/lib/html"
)

var udb authlib.IUsers

func init() {
	udb = authlib.NewUsers()
}

//Router to include auth functions in web site
func Router() http.Handler {
	r := html.NewRouter("auth")
	r.Form("register", register)
	r.Form("activate", activate)
	r.Form("reset", reset)
	r.Form("login", login)
	return r
}

func register(res html.ISessionRes, req *http.Request) (string, error) {
	req.ParseForm()
	u, err := udb.Add(req.FormValue("email"))
	if err != nil {
		return "", log.Wrapf(err, "Cannot register that email address")
	}

	//registered, send activation email
	log.Debugf("Activation with email: %s and tpw: %s", u.Email(), u.Tpw())

	//todo: send this in email...
	urlValues := url.Values{}
	urlValues.Set("email", u.Email())
	urlValues.Set("tpw", u.Tpw())
	activationLink := fmt.Sprintf("/auth/activate?%s", urlValues.Encode())

	//show confirmation page
	return res.Message("Success",
		"Your registration was successful.<br>"+
			"An email was sent to your address %s<br>"+
			"TO REMOVE:  <a href=\"%s\">Activate Here</a>"+
			"Follow the instructions in the email message to activate your account.",
		u.Email(),
		activationLink)
}

func activate(res html.ISessionRes, req *http.Request) (string, error) {
	req.ParseForm()
	log.Debugf("Activate Form Values: %+v", req.Form)
	u := udb.GetByEmail(req.FormValue("email"))
	if u == nil {
		return "/", log.Wrapf(nil, "Your email address is not yet registered")
	}

	tpw := req.FormValue("tpw")
	if len(tpw) != 8 {
		return "/", log.Wrapf(nil, "Invalid link")
	}

	rpw := req.FormValue("rpw")
	if len(rpw) == 0 {
		return "", log.Wrapf(nil, "New password not specified")
	}

	cpw := req.FormValue("cpw")
	if len(rpw) == 0 || rpw != cpw {
		return "", log.Wrapf(nil, "New password not repeated")
	}

	if err := u.SetPassword(tpw, rpw); err != nil {
		return "", log.Wrapf(err, "Failed to change your password. Please try again later")
	}

	return res.Message("Success", "Your password was changed as specified.<br>"+
		" You may now <a href=\"/auth/login?email=%s\">login</a>.", u.Email())
}

func reset(res html.ISessionRes, req *http.Request) (string, error) {
	req.ParseForm()
	u := udb.GetByEmail(req.FormValue("email"))
	if u == nil {
		return "/", log.Wrapf(nil, "That email address is not yet registered.")
	}

	if err := u.Reset(); err != nil {
		return res.Message("Error", "Failed to reset password")
	}

	//Todo: Send email...
	//registered, send activation email
	log.Debugf("Activation with email: %s and tpw: %s", u.Email(), u.Tpw())

	//show confirmation page
	return res.Message("Success",
		"Your password was successfully reset.<br>"+
			"An email was sent to your address %s<br>"+
			"Follow the instructions in the email message to change your password.", u.Email())
}

func login(res html.ISessionRes, req *http.Request) (string, error) {
	req.ParseForm()
	u := udb.GetByEmail(req.FormValue("email"))
	if u == nil {
		return "", log.Wrapf(nil, "Your email address is not yet registered.")
	}

	if ok := u.Auth(req.FormValue("rpw")); !ok {
		return "", log.Wrapf(nil, "Cannot authenticate you with that email and password")
	}

	//todo - show user's home page and manage the session...

	res.Session().SetUser(u)
	return "/user/home", nil
}

func logout(res html.ISessionRes, req *http.Request) (string, error) {
	res.Session().SetUser(nil)
	return "/", nil
}
