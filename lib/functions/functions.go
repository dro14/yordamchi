package functions

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dro14/yordamchi/lib/types"
	"github.com/gin-gonic/gin"
)

func LanguageCode(lang string) string {
	if lang == "" {
		lang = "uz"
	} else if lang != "uz" && lang != "ru" {
		lang = "en"
	}
	return lang
}

func Sleep(retryDelay *time.Duration) {
	log.Printf("retrying request after %v", *retryDelay)
	time.Sleep(*retryDelay)
	*retryDelay *= 2
}

func Transaction(params *types.Params) string {
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%s_%d_%d", params.Account.OrderID, params.Time, params.Amount)))
	hashValue := hash.Sum(nil)
	return hex.EncodeToString(hashValue)
}

var MerchantKey string

func Authorized(c *gin.Context) bool {

	header := c.GetHeader("Authorization")
	header = strings.TrimPrefix(header, "Basic ")
	reader := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(header))

	bts := make([]byte, 50)
	n, err := reader.Read(bts)
	if err != nil {
		log.Printf("can't read header: %v", err)
		c.JSON(200, gin.H{
			"error": gin.H{
				"code":    -32504,
				"message": "Unauthorized",
			},
		})
		return false
	}

	header = string(bts[:n])
	if header != "Paycom:"+MerchantKey {
		log.Printf("unauthorized: %s", header)
		c.JSON(200, gin.H{
			"error": gin.H{
				"code":    -32504,
				"message": "Unauthorized",
			},
		})
		return false
	}

	return true
}

//func MarkdownV2(text string) string {
//
//	var (
//		found       bool
//		before      []byte
//		buffer      bytes.Buffer
//		char        = []byte{0}
//		bts         = []byte(text)
//		escape      = []byte{'\\', 0}
//		block       = []byte{'`', '`', '`'}
//		blocks      = bytes.Count(bts, block)
//		escapeChars = []byte{'\\', '_', '*', '[', ']', '(', ')', '~', '>', '#', '+', '-', '=', '|', '{', '}', '.', '!', '`'}
//	)
//
//	for i := range escapeChars {
//		char[0] = escapeChars[i]
//		escape[1] = escapeChars[i]
//		if i != 18 {
//			bts = bytes.ReplaceAll(bts, char, escape)
//		}
//	}
//
//	if blocks%2 != 0 {
//		bts = append(bts, block...)
//	}
//
//	for {
//		before, bts, found = bytes.Cut(bts, block)
//		if bytes.Count(before, char)%2 != 0 {
//			before = bytes.ReplaceAll(before, char, escape)
//		}
//		buffer.Write(before)
//		if !found {
//			break
//		}
//		buffer.Write(block)
//
//		before, bts, _ = bytes.Cut(bts, block)
//		before = bytes.ReplaceAll(before, char, escape)
//		buffer.Write(before)
//		buffer.Write(block)
//	}
//
//	return buffer.String()
//}
