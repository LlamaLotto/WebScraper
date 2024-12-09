package Model

type NLBLottery struct {
	Date       string   `json:"date"`
	DrawNumber string   `json:"draw_number"`
	Numbers    []string `json:"numbers"`
}
