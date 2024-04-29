package other

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/dro14/yordamchi/utils"
)

func (o *APIs) OCR(ctx context.Context, photoURL, caption string) string {
	var buffer bytes.Buffer
	request := map[string]string{"url": photoURL}
	err := json.NewEncoder(&buffer).Encode(request)
	if err != nil {
		log.Println("can't encode request:", err)
		return caption
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, o.ocrURL, &buffer)
	if err != nil {
		log.Println("can't create request:", err)
		return caption
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Ocp-Apim-Subscription-Key", o.subscriptionKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("can't do request:", err)
		return caption
	}
	defer func() { _ = resp.Body.Close() }()

	bts, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("can't read body:", err)
		return caption
	}

	response := Response{}
	err = json.Unmarshal(bts, &response)
	if err != nil {
		log.Println("can't decode response:", err)
		return caption
	}

	if caption != "" {
		caption = "\n\n" + caption
	}
	return response.ReadResult.Content + caption
}

func (o *APIs) Translate(sl, tl, q string) string {
	parts := strings.Split(q, "```")
	for i := 0; i < len(parts); i += 2 {
		qs := utils.Slice(parts[i], 5000)

		for j := range qs {
			resp, err := http.Get(fmt.Sprintf(o.translateURL, sl, tl, url.QueryEscape(qs[j])))
			if err != nil {
				log.Println("can't do request:", err)
				return strings.Join(qs, " ")
			}

			if resp.StatusCode != http.StatusOK {
				log.Println("bad status:", resp.Status)
				return strings.Join(qs, " ")
			}

			bts, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Println("can't read body:", err)
				return strings.Join(qs, " ")
			}

			_, after, found := strings.Cut(string(bts), "<div class=\"result-container\">")
			if !found {
				log.Println("not found")
				return strings.Join(qs, " ")
			}

			translation, _, found := strings.Cut(after, "</div>")
			if !found {
				log.Println("not found")
				return strings.Join(qs, " ")
			}

			qs[j] = html.UnescapeString(translation)
		}
		parts[i] = strings.Join(qs, " ")
	}

	return strings.Join(parts, "```")
}
