package structure

import "time"

type Model struct {
	OrderUid          string    `json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Delivery          Delivery  `json:"delivery"`
	Payment           Payment   `json:"payment"`
	Items             []Item    `json:"items"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerId        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	ShardKey          string    `json:"shardkey"`
	SmId              int       `json:"sm_id" fake:"{number:1,100}"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transcation  string `json:"transaction"`
	RequestId    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount" fake:"{number:1,100}"`
	PaymentDt    int    `json:"payment_dt" fake:"{number:1,100}"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost" fake:"{number:1,100}"`
	GoodsTotal   int    `json:"goods_total" fake:"{number:1,100}"`
	CustomFee    int    `json:"custom_fee" fake:"{number:1,100}"`
}

type Item struct {
	ChrtId      int    `json:"chrt_id" fake:"{number:1,100}"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price" fake:"{number:1,100}"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale" fake:"{number:1,10000}"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price" fake:"{number:1,100000}"`
	NmId        int    `json:"nm_id" fake:"{number:1,100}"`
	Brand       string `json:"brand"`
	Status      int    `json:"status" fake:"{number:1,100}"`
}
