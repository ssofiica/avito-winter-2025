package entity

type Merch struct {
	ID   uint32
	Name string
	Cost uint32
}

type InfoResponse struct {
	Coins       uint32      `json:"coins"`
	Inventory   []Inventory `json:"inventory"`
	CoinHistory CoinHistory `json:"coinHistory"`
}

type Inventory struct {
	Type     string `json:"type"`
	Quantity uint32 `json:"quantity"`
}

type CoinHistory struct {
	Received []Received `json:"received"`
	Sent     []Sent     `json:"sent"`
}

type Received struct {
	FromUser string `json:"fromUser"`
	Amount   uint32 `json:"amount"`
}

type Sent struct {
	ToUser string `json:"toUser"`
	Amount uint32 `json:"amount"`
}
