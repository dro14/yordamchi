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
	postgres   *postgres.Postgres
}

func New() *Payme {
	merchantID, ok := os.LookupEnv("MERCHANT_ID")
	if !ok {
		log.Fatal("merchant id is not specified")
	}

	testKey, ok := os.LookupEnv("TEST_KEY")
	if !ok {
		log.Fatal("test key is not specified")
	}

	realKey, ok := os.LookupEnv("REAL_KEY")
	if !ok {
		log.Fatal("real key is not specified")
	}

	return &Payme{
		merchantID: merchantID,
		testKey:    testKey,
		realKey:    realKey,
		postgres:   postgres.New(),
	}
}
