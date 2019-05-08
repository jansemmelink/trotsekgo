package email

import (
	"regexp"

	"github.com/jansemmelink/log"
)

var emailAddressPattern *regexp.Regexp

func init() {
	id := `[a-zA-Z0-9][a-zA-Z0-9_-]*`
	ids := id + `(\.` + id + `)*`
	emailAddressPattern = regexp.MustCompile(`^` + ids + `@` + id + `\.` + ids + `$`)
}

//CheckAddress validate email address format
func CheckAddress(email string) error {
	if !emailAddressPattern.MatchString(email) {
		return log.Wrapf(nil, "Invalid email address")
	}
	return nil
}
