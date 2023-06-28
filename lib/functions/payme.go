package functions

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	"github.com/dro14/yordamchi/lib/constants"
	"github.com/dro14/yordamchi/lib/types"
	"github.com/gin-gonic/gin"
)

func Transaction(params *types.Params) string {
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%s_%d_%d", params.Account.OrderID, params.Time, params.Amount)))
	hashValue := hash.Sum(nil)
	return hex.EncodeToString(hashValue)
}

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
	if header != "Paycom:"+constants.TestKey && header != "Paycom:"+constants.RealKey {
		log.Printf("unauthorized header: %s", header)
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
