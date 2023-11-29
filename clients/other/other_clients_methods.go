package other

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dro14/yordamchi/utils"
	"html"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func (o *APIs) Vision(ctx context.Context, photoURL, caption string) string {
	var buffer bytes.Buffer
	request := map[string]string{"url": photoURL}
	err := json.NewEncoder(&buffer).Encode(request)
	if err != nil {
		log.Println("can't encode request:", err)
		return caption
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, o.visionURL, &buffer)
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

	response := Response{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Println("can't decode response:", err)
		return caption
	}

	var builder strings.Builder
	if len(caption) > 0 {
		builder.WriteString(caption + ":\n\n")
	}

	for _, region := range response.Regions {
		for _, line := range region.Lines {
			for _, word := range line.Words {
				builder.WriteString(word.Text + " ")
			}
			builder.WriteString("\n")
		}
		builder.WriteString("\n\n")
	}

	return builder.String()
}

func (o *APIs) Translate(sl, tl, q string) string {
	qs := utils.Slice(q, 5000)

	for i := range qs {
		resp, err := http.Get(fmt.Sprintf(o.translateURL, sl, tl, url.QueryEscape(qs[i])))
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

		qs[i] = html.UnescapeString(translation)
	}

	return strings.Join(qs, " ")
}

func (o *APIs) Search(ctx context.Context, query string) string {
	var buffer bytes.Buffer
	request := map[string]any{
		"query":   query,
		"lang":    ctx.Value("language_code").(string),
		"user_id": ctx.Value("user_id").(int64),
	}
	err := json.NewEncoder(&buffer).Encode(request)
	if err != nil {
		log.Println("can't encode request:", err)
		return ""
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, o.searchURL, &buffer)
	if err != nil {
		log.Println("can't create request:", err)
		return ""
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("can't do request:", err)
		return ""
	}
	defer func() { _ = resp.Body.Close() }()

	response := make(map[string]any)
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Println("can't decode response:", err)
		return ""
	}

	if response["success"] == false {
		return ""
	}

	return response["results"].(string)
}
