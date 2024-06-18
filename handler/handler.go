package handler

import (
	"log"

	"github.com/dro14/yordamchi/legacy"
	"github.com/dro14/yordamchi/payment/click"
	"github.com/dro14/yordamchi/payment/click/methods"
	clickTypes "github.com/dro14/yordamchi/payment/click/types"
	"github.com/dro14/yordamchi/payment/payme"
	paymeTypes "github.com/dro14/yordamchi/payment/payme/types"
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
	click     *click.Click
}

func New() *Handler {
	return &Handler{
		router:    gin.Default(),
		processor: processor.New(),
		legacy:    legacy.New(),
		payme:     payme.New(),
		click:     click.New(),
	}
}

func (h *Handler) Run(port string) error {
	h.router.GET("/", h.Root)
	h.router.GET("/logs", h.Logs)
	h.router.POST("/main", h.Main)
	h.router.POST("/legacy", h.Legacy)
	h.router.POST("/payme", h.Payme)
	h.router.POST("/click/prepare", h.Click)
	h.router.POST("/click/complete", h.Click)
	return h.router.Run(":" + port)
}

func (h *Handler) Root(c *gin.Context) {
	c.JSON(200, gin.H{"ok": true})
}

func (h *Handler) Logs(c *gin.Context) {
	utils.SendLogFiles()
	c.JSON(200, gin.H{"ok": true})
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
	request, response := &paymeTypes.Request{}, gin.H{}
	err := c.ShouldBindJSON(request)
	switch {
	case err != nil:
		log.Println("can't bind Payme body:", err)
		response = gin.H{"error": gin.H{"code": -32700, "message": "Parse error"}}
	default:
		response = h.payme.Process(c, request)
	}
	c.JSON(200, response)
}

func (h *Handler) Click(c *gin.Context) {
	request, response := &clickTypes.Request{}, gin.H{}
	err := c.ShouldBind(request)
	switch {
	case err != nil:
		log.Println("can't bind Click body:", err)
		response = gin.H{"error": -8, "error_note": "Error in request from click"}
	case request.SignString != h.click.SingString(request):
		response = gin.H{"error": -1, "error_note": "SIGN CHECK FAILED!"}
	case request.Error != 0:
		response = h.click.Cancel(request)
	case request.Action == methods.Prepare:
		response = h.click.Prepare(request)
	case request.Action == methods.Complete:
		response = h.click.Complete(request)
	default:
		response = gin.H{"error": -3, "error_note": "Action not found"}
	}
	c.JSON(200, response)
}
