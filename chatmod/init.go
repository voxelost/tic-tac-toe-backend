package chatmod

import (
	"io"
	"net/http"
	"strings"
)

const LIST_URL = "https://raw.githubusercontent.com/LDNOOBW/List-of-Dirty-Naughty-Obscene-and-Otherwise-Bad-Words/master/en"

var badWordsCache []string

func init() {
	resp, err := http.Get(LIST_URL)
	if err != nil {
		return
	}
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	badWordsCache = strings.Split(string(respBytes), "\n")
}
