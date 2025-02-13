package entity

type SendCoinRequest struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

func (s SendCoinRequest) Valid() bool {
	return s.ToUser != "" && s.Amount > 0
}

type Transaction struct {
	From   uint32
	To     uint32
	Amount uint32
}
