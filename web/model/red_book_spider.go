package model

type RedBookSpider struct {
	Id          int    `json:"id"`
	Link        string `json:"link"`
	CreatedAt   int    `json:"created_at"`
	SaveJson    string `json:"save_json"`
	SuccessRate string `json:"success_rate"`
}

var RedBookSpiderModel RedBookSpider

func (RedBookSpider) Update(log *RedBookSpider) {

}
