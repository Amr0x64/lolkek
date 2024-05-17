package models

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
)

type Delivery struct {
	Name    string `json:"name" validate:"required,max=30"`
	Phone   string `json:"phone" validate:"required,max=20"`
	Zip     string `json:"zip" validate:"required,max=10"`
	City    string `json:"city" validate:"required,max=50"`
	Address string `json:"address" validate:"required,max=50"`
	Region  string `json:"region" validate:"required,max=50"`
	Email   string `json:"email" validate:"required,email"`
}

type Item struct {
	ChrtId      int    `json:"chrt_id" validate:"required"`
	TrackNumber string `json:"track_number" validate:"required,min=14,max=14"`
	Price       int    `json:"price" validate:"gt=0"`
	Rid         string `json:"rid" validate:"required,min=21,max=21"`
	Name        string `json:"name" validate:"required"`
	Sale        int    `json:"sale" validate:"required"`
	Size        string `json:"size" validate:"required"`
	TotalPrice  int    `json:"total_price" validate:"gt=0"`
	NmId        int    `json:"nm_id" validate:"required"`
	Brand       string `json:"brand" validate:"required"`
	Status      int    `json:"status" validate:"required"`
	OrderUid    string `json:"-" db:"order_uid"`
}

type Payment struct {
	Transaction  string `json:"transaction" validate:"required"`
	RequestId    string `json:"request_id"`
	Currency     string `json:"currency" validate:"required"`
	Provider     string `json:"provider" validate:"required"`
	Amount       int    `json:"amount" validate:"gt=0"`
	PaymentDt    int    `json:"payment_dt" validate:"required"`
	Bank         string `json:"bank" validate:"required"`
	DeliveryCost int    `json:"delivery_cost" validate:"gt=0"`
	GoodsTotal   int    `json:"goods_total" validate:"gt=0"`
	CustomFee    int    `json:"custom_fee"`
}

type Order struct {
	OrderUid          string    `json:"order_uid" validate:"required,min=19,max=19"`
	TrackNumber       string    `json:"track_number" validate:"required,min=14,max=14"`
	Entry             string    `json:"entry" validate:"required"`
	Locale            string    `json:"locale" validate:"required,oneof=ru en"`
	InternalSignature string    `json:"internal_signature"`
	CustomerId        string    `json:"customer_id" validate:"required"`
	DeliveryService   string    `json:"delivery_service" validate:"required"`
	ShardKey          string    `json:"shardkey" validate:"required"`
	SmId              int       `json:"sm_id" validate:"required"`
	DateCreated       time.Time `json:"date_created" format:"2006-01-02T06:22:19Z" validate:"required"`
	OofShard          string    `json:"oof_shard" validate:"required"`
	Delivery          Delivery  `json:"delivery" validate:"required"`
	Payment           Payment   `json:"payment" validate:"required"`
	Items             []Item    `json:"items" validate:"required"`
}

func (o *Order) Validate() error {
	validate := validator.New()

	if err := validate.Struct(o); err != nil {
		return err
	}

	// Валидация вложенных структур
	if err := validate.Struct(o.Delivery); err != nil {
		return err
	}
	if err := validate.Struct(o.Payment); err != nil {
		return err
	}
	for _, item := range o.Items {
		if err := validate.Struct(item); err != nil {
			return err
		}
	}

	return nil
}

func (i *Item) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		log.Println("failed to unmarshal JSON value")
		return errors.New("1")
	}
	return json.Unmarshal(bytes, i)
}
