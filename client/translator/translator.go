package translator

import (
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/dro14/yordamchi/lib/constants"
	"github.com/dro14/yordamchi/lib/functions"
)

func Translate(sl, tl, q string) (string, error) {

	qs := functions.Slice(q)

	for i := range qs {
		retryDelay := constants.RetryDelay
		attempts := 0
	Retry:
		attempts++
		resp, err := http.Get(fmt.Sprintf("https://translate.google.com/m?sl=%s&tl=%s&q=%s", sl, tl, url.QueryEscape(qs[i])))
		if err != nil {
			log.Printf("can't do request: %v", err)
			return "", err
		}

		if resp.StatusCode != http.StatusOK {
			log.Printf("bad status: %s", resp.Status)
			if attempts < constants.RetryAttempts {
				functions.Sleep(&retryDelay)
				goto Retry
			}
			return "", fmt.Errorf("bad status: %s", resp.Status)
		}

		bts, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("can't read body: %v", err)
			return "", err
		}

		_, after, found := strings.Cut(string(bts), "<div class=\"result-container\">")
		if !found {
			log.Printf("not found")
			return "", fmt.Errorf("not found")
		}

		translation, _, found := strings.Cut(after, "</div>")
		if !found {
			log.Printf("not found")
			return "", fmt.Errorf("not found")
		}

		qs[i] = html.UnescapeString(translation)
	}

	return strings.Join(qs, " "), nil
}
