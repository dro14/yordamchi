package ocr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"net/http"
)

type ComputerVisionRequest struct {
	URL string `json:"url"`
}

type ComputerVisionResponse struct {
	Regions []struct {
		Lines []struct {
			Words []struct {
				Text string `json:"text"`
			} `json:"words"`
		} `json:"lines"`
	} `json:"regions"`
}

func main() {
	subscriptionKey := "<your-subscription-key>"
	endpoint := "<your-endpoint>"

	imageURL := "https://example.com/path/to/your/image.jpg"
	visionRequest := &ComputerVisionRequest{URL: imageURL}
	jsonRequest, err := json.Marshal(visionRequest)
	if err != nil {
		panic(err)
	}

	url := endpoint + "/vision/v3.1/ocr"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonRequest))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Ocp-Apim-Subscription-Key", subscriptionKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var visionResponse ComputerVisionResponse
	err = json.Unmarshal(body, &visionResponse)
	if err != nil {
		panic(err)
	}

	for _, region := range visionResponse.Regions {
		for _, line := range region.Lines {
			for _, word := range line.Words {
				fmt.Print(word.Text + " ")
			}
			fmt.Println()
		}
	}
}

func Handler(c *gin.Context) {

	update := &tgbotapi.Update{}
	err := c.ShouldBindJSON(update)
	if err != nil {
		log.Printf("can't bind json: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"ok": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
	if update.Message == nil || update.Message.Photo == nil {
		return
	}

	photo := *update.Message.Photo
	photoFileID := photo[len(photo)-1].FileID

	fmt.Println(photoFileID)
}
