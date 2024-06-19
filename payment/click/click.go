package click

import (
	"log"
	"os"
	"strconv"

	"github.com/dro14/yordamchi/storage/postgres"
)

type Click struct {
	merchantID int
	serviceID  int
	secretKey  string
	url        string
	postgres   *postgres.Postgres
}

func New() *Click {
	value, ok := os.LookupEnv("CLICK_MERCHANT_ID")
	if !ok {
		log.Fatal("Click merchant id is not specified")
	}
	merchantID, err := strconv.Atoi(value)
	if err != nil {
		log.Fatal("expected an integer for Click merchant id, but got:", value)
	}

	value, ok = os.LookupEnv("CLICK_SERVICE_ID")
	if !ok {
		log.Fatal("Click service id is not specified")
	}
	serviceID, err := strconv.Atoi(value)
	if err != nil {
		log.Fatal("expected an integer for Click service id, but got:", value)
	}

	secretKey, ok := os.LookupEnv("CLICK_SECRET_KEY")
	if !ok {
		log.Fatal("Click secret key is not specified")
	}

	return &Click{
		merchantID: merchantID,
		serviceID:  serviceID,
		secretKey:  secretKey,
		url:        "https://my.click.uz/services/pay?service_id=%d&merchant_id=%d&amount=%d&transaction_param=%d&return_url=https://t.me/yordamchi_ai_bot",
		postgres:   postgres.New(),
	}
}
