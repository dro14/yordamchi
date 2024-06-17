package payme

import (
	"log"
	"os"

	"github.com/dro14/yordamchi/storage/postgres"
)

type Payme struct {
	merchantID string
	testKey    string
	realKey    string
	url        string
	postgres   *postgres.Postgres
}

func New() *Payme {
	merchantID, ok := os.LookupEnv("PAYME_MERCHANT_ID")
	if !ok {
		log.Fatal("Payme merchant id is not specified")
	}

	testKey, ok := os.LookupEnv("PAYME_TEST_KEY")
	if !ok {
		log.Fatal("Payme test key is not specified")
	}

	realKey, ok := os.LookupEnv("PAYME_REAL_KEY")
	if !ok {
		log.Fatal("Payme real key is not specified")
	}

	return &Payme{
		merchantID: merchantID,
		testKey:    testKey,
		realKey:    realKey,
		url:        "https://checkout.paycom.uz/",
		postgres:   postgres.New(),
	}
}
