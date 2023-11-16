package payme

import (
	"bytes"
	"encoding/base64"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

func (p *Payme) authorized(c *gin.Context) gin.H {
	unauthorized := gin.H{"error": gin.H{"code": -32504, "message": "Unauthorized"}}
	header := c.GetHeader("Authorization")
	header = strings.TrimPrefix(header, "Basic ")
	reader := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(header))
	bts := make([]byte, 50)

	n, err := reader.Read(bts)
	if err != nil {
		log.Println("can't read header:", err)
		return unauthorized
	}
	header = string(bts[:n])
	if header != "Paycom:"+p.testKey && header != "Paycom:"+p.realKey {
		log.Println("unauthorized header:", header)
		return unauthorized
	}
	return nil
}
