package other

type Response struct {
	ReadResult struct {
		Content string `json:"content"`
	} `json:"readResult"`
}
