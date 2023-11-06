package google

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func GoogleSearch(query string) string {

	request := make(map[string]string)
	request["query"] = query

	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(request)
	if err != nil {
		log.Printf("can't encode request: %v", err)
		return ""
	}

	resp, err := http.Post("https://google.victoriousriver-fffd2d70.westeurope.azurecontainerapps.io/search", "application/json", &buffer)
	if err != nil {
		log.Printf("can't send request: %v", err)
		return ""
	}

	response := make(map[string]string)
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return ""
	}
	_ = resp.Body.Close()

	return response["results"]
}
