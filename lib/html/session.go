package html

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/jansemmelink/log"
	"github.com/jansemmelink/trotsek/lib/auth"
)

//NewSession ...
func NewSession() ISession {
	return &session{}
}

//ISession ...
type ISession interface {
	ID() string
	User() auth.IUser     //nil when not logged in
	SetUser(u auth.IUser) //set to nil when logging out
	Res(res http.ResponseWriter) ISessionRes
}

type session struct {
	id   string
	user auth.IUser
}

func (s session) ID() string {
	return s.id
}

func (s session) User() auth.IUser {
	return s.user
}

func (s *session) SetUser(u auth.IUser) {
	s.user = u
}

func (s *session) Res(res http.ResponseWriter) ISessionRes {
	return sessionRes{
		ResponseWriter: res,
		session:        s,
	}
}

//ISessionRes wraps the response writer to provide consistent site responses
type ISessionRes interface {
	http.ResponseWriter
	Session() ISession
	Message(title, format string, args ...interface{}) (string, error)
	FormTemplate(filename string, data map[string]string)
}

type sessionRes struct {
	http.ResponseWriter //embedded
	session             ISession
}

func (sr sessionRes) Write(data []byte) (int, error) {
	//implement page header/footer here... or prevent this from being called!
	return sr.ResponseWriter.Write(data)
}

func (sr sessionRes) Session() ISession {
	return sr.session
}

//Message writes a message page
func (sr sessionRes) Message(title, format string, args ...interface{}) (string, error) {
	sr.pageHeader(sr.ResponseWriter)

	//todo: use safe html encoding not to break sites with controls in text...
	messageHTML := fmt.Sprintf("<h1>%s</h1><p>%s</p>", title, fmt.Sprintf(format, args...))
	sr.ResponseWriter.Write([]byte(messageHTML))

	sr.pageFooter(sr.ResponseWriter)
	return "", nil
}

func (sr sessionRes) FormTemplate(filename string, data map[string]string) {
	log.Debugf("FormTemplate(%s) data=%+v ...", filename, data)
	t, err := template.ParseFiles(filename)
	if err != nil {
		sr.Message("Error", "Page Not Found")
		panic(log.Wrapf(err, "Failed to load template: %s"+filename))
	}

	sr.pageHeader(sr.ResponseWriter)
	err = t.Execute(sr.ResponseWriter, data)
	if err != nil {
		panic(log.Wrapf(err, "failed to execute template %s", filename))
	}
	sr.pageFooter(sr.ResponseWriter)
}

const headerTemplate = `
<HTML>
	<HEAD>
		<link rel="stylesheet" href="/resources/page-style.css">
		<link rel="stylesheet" href="/resources/menu-style.css">
		<link rel="stylesheet" href="/resources/trotsek.css">
	</HEAD>
	<BODY>
		<DIV id="nav">
			<ul>
				<li><a href="/">Home</a></li>
				<li><a href="/wallet">Wallet</a></li>
				<li><a href="/events">Events</a></li>
				<li><a href="/entries">Entries</a></li>
				<li><a href="/profile">Profile</a></li>
				<li><a href="/auth/register">Register</a></li>
				<li><a href="/auth/login">Login</a></li>
			</ul>
		</DIV>
		<DIV class="contents">
			<!-- END OF PAGE HEADER -->
`

const footerTemplate = `
			<!-- START OF PAGE FOOTER -->
		</DIV>
	</BODY>
</HTML>
`

func (sr sessionRes) pageHeader(res http.ResponseWriter) {
	t, err := template.New("header").Parse(headerTemplate)
	if err != nil {
		log.Errorf("Failed to render page header: %v", err)
		panic(err)
	}

	data := []string{}
	err = t.Execute(res, data)
	if err != nil {
		panic(log.Wrapf(err, "failed to execute template"))
	}
}

func (sr sessionRes) pageFooter(res http.ResponseWriter) {
	t, err := template.New("footer").Parse(footerTemplate)
	if err != nil {
		log.Errorf("Failed to render page footer: %v", err)
		panic(err)
	}

	data := []string{}
	err = t.Execute(res, data)
	if err != nil {
		panic(log.Wrapf(err, "failed to execute template"))
	}
}
