package data

import (
	"html"
	"regexp"
)

func SafeString(fileName string) string {
	re := regexp.MustCompile(`[<>'"$\{\}]`)
	return re.ReplaceAllString(html.EscapeString(fileName), "")
}
