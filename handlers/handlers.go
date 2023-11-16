package handlers

import (
	"log"

	"github.com/dro14/yordamchi/legacy"
	"github.com/dro14/yordamchi/payme"
	"github.com/dro14/yordamchi/payme/types"
	"github.com/dro14/yordamchi/processor"
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
	h.router.POST("/main", h.Main)
	h.router.POST("/legacy", h.Legacy)
	h.router.POST("/payme", h.Payme)
	return h.router.Run(":" + port)
}

func (h *Handler) Main(c *gin.Context) {
	update := &tgbotapi.Update{}
	err := c.ShouldBindJSON(update)
	if err != nil {
		log.Println("can't bind json:", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	switch {
	case update.Message != nil:
		go h.processor.Message(update.Message)
	case update.CallbackQuery != nil:
		go h.processor.CallbackQuery(update.CallbackQuery)
	case update.MyChatMember != nil:
		go h.processor.MyChatMember(update.MyChatMember)
	default:
		log.Printf("unknown update type:\n%+v", update)
	}
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
	switch {
	case update.Message != nil:
		h.legacy.Redirect(update.Message)
	default:
		log.Printf("unknown update type:\n%+v", update)
	}
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
