package types

type Request struct {
	ClickTransID      int64   `json:"click_trans_id"`
	ServiceID         int     `json:"service_id"`
	ClickPaydocID     int64   `json:"click_paydoc_id"`
	MerchantTransID   string  `json:"merchant_trans_id"`
	MerchantPrepareID int     `json:"merchant_prepare_id"`
	Amount            float64 `json:"amount"`
	Action            int     `json:"action"`
	Error             int     `json:"error"`
	ErrorNote         string  `json:"error_note"`
	SignTime          int64   `json:"sign_time"`
	SignString        string  `json:"sign_string"`
}

type Response struct {
	ClickTransID      int64  `json:"click_trans_id"`
	MerchantTransID   string `json:"merchant_trans_id"`
	MerchantPrepareID int    `json:"merchant_prepare_id,omitempty"`
	MerchantConfirmID int    `json:"merchant_confirm_id,omitempty"`
	Error             int    `json:"error"`
	ErrorNote         string `json:"error_note"`
}
