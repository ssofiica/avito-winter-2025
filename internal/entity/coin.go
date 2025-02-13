package entity

type SendCoinRequest struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

func (s SendCoinRequest) Valid() bool {
	return s.ToUser != "" && s.Amount > 0
}
