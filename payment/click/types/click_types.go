package types

type Request struct {
	ClickTransID      int64   `form:"click_trans_id" binding:"required"`
	ServiceID         int     `form:"service_id" binding:"required"`
	ClickPaydocID     int64   `form:"click_paydoc_id" binding:"required"`
	MerchantTransID   int     `form:"merchant_trans_id" binding:"required"`
	MerchantPrepareID int     `form:"merchant_prepare_id"`
	Amount            float64 `form:"amount" binding:"required"`
	Action            int     `form:"action"`
	Error             int     `form:"error"`
	ErrorNote         string  `form:"error_note" binding:"required"`
	SignTime          string  `form:"sign_time" binding:"required"`
	SignString        string  `form:"sign_string" binding:"required"`
}
