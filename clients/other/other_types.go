package other

type Response struct {
	Regions []struct {
		Lines []struct {
			Words []struct {
				Text string `json:"text"`
			} `json:"words"`
		} `json:"lines"`
	} `json:"regions"`
}
