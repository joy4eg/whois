package whois

import "regexp"

var tldRex = regexp.MustCompile(`^\.(xn--)?[a-z0-9]+$`)

func matchesTLD(s string) bool {
	return tldRex.MatchString(s)
}
