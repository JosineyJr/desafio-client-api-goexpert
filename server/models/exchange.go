package models

type Exchange struct {
	Code        string
	Codein      string
	Name        string
	High        string
	Low         string
	VarBid      string
	PctChange   string
	Bid         string
	Ask         string
	Timestamp   string
	Create_date string
}

type IExchangeRepository interface {
	Insert(exchange Exchange) error
}
