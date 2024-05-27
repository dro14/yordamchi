package handler

import (
	"log"

	"github.com/dro14/yordamchi/legacy"
	"github.com/dro14/yordamchi/payme"
	"github.com/dro14/yordamchi/payme/types"
	"github.com/dro14/yordamchi/processor"
	"github.com/dro14/yordamchi/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	router    *gin.Engine
	processor *processor.Processor
	legacy    *legacy.Legacy
	payme     *payme.Payme
}

func New() *Handler {
	return &Handler{
		router:    gin.Default(),
		processor: processor.New(),
		legacy:    legacy.New(),
		payme:     payme.New(),
	}
}

func (h *Handler) Run(port string) error {
	h.router.GET("/", h.Root)
	h.router.POST("/main", h.Main)
	h.router.POST("/legacy", h.Legacy)
	h.router.POST("/payme", h.Payme)
	h.router.GET("/logs", h.Logs)
	return h.router.Run(":" + port)
}

func (h *Handler) Root(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Hello, World!"})
}

func (h *Handler) Main(c *gin.Context) {
	update := &tgbotapi.Update{}
	err := c.ShouldBindJSON(update)
	if err != nil {
		log.Println("can't bind json:", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	go h.processor.Update(update)
	c.JSON(200, gin.H{"ok": true})
}

func (h *Handler) Legacy(c *gin.Context) {
	update := &tgbotapi.Update{}
	err := c.ShouldBindJSON(update)
	if err != nil {
		log.Println("can't bind json:", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	h.legacy.Redirect(update)
	c.JSON(200, gin.H{"ok": true})
}

func (h *Handler) Payme(c *gin.Context) {
	request := &types.Request{}
	err := c.ShouldBindJSON(request)
	if err != nil {
		log.Println("can't bind json:", err)
		c.JSON(200, gin.H{"error": gin.H{"code": -32700, "message": "Parse error"}})
		return
	}
	response := h.payme.Respond(c, request)
	c.JSON(200, response)
}

func (h *Handler) Logs(c *gin.Context) {
	utils.SendLogFiles()
	c.JSON(200, gin.H{"ok": true})
}
