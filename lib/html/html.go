package html

import (
	"fmt"
	"net/http"
)

//Message writes a message page
func Message(res http.ResponseWriter, title, format string, args ...interface{}) {
	//todo properly with header/footer from website!
	res.Write([]byte(fmt.Sprintf("<h1>%s</h1>", title)))
	res.Write([]byte(fmt.Sprintf("<p>%s</p>", fmt.Sprintf(format, args...))))
	//todo: html encoding not to break sites with controls in text...
}
