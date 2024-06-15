package types

type Request struct {
	ClickTransID      int64   `form:"click_trans_id"`
	ServiceID         int     `form:"service_id"`
	ClickPaydocID     int64   `form:"click_paydoc_id"`
	MerchantTransID   string  `form:"merchant_trans_id"`
	MerchantPrepareID int     `form:"merchant_prepare_id"`
	Amount            float64 `form:"amount"`
	Action            int     `form:"action"`
	Error             int     `form:"error"`
	ErrorNote         string  `form:"error_note"`
	SignTime          string  `form:"sign_time"`
	SignString        string  `form:"sign_string"`
}

type Response struct {
	ClickTransID      int64  `json:"click_trans_id"`
	MerchantTransID   string `json:"merchant_trans_id"`
	MerchantPrepareID int    `json:"merchant_prepare_id,omitempty"`
	MerchantConfirmID int    `json:"merchant_confirm_id,omitempty"`
	Error             int    `json:"error"`
	ErrorNote         string `json:"error_note"`
}
