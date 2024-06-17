package payme

import (
	"bytes"
	"context"
	"encoding/base64"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

func id(ctx context.Context) int64 {
	return ctx.Value("user_id").(int64)
}

func lang(ctx context.Context) string {
	return ctx.Value("language_code").(string)
}

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
