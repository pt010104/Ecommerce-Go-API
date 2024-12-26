package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID            primitive.ObjectID `bson:"_id"`
	UserID        primitive.ObjectID `bson:"user_id"`
	CheckoutID    primitive.ObjectID `bson:"checkout_id"`
	ShopID        primitive.ObjectID `bson:"shop_id"`
	Products      []OrderProduct     `bson:"products"`
	Status        string             `bson:"status"`
	TotalPrice    float64            `bson:"total_price"`
	Address       Address            `bson:"address"`
	PaymentMethod string             `bson:"payment_method"`
	VoucherID     primitive.ObjectID `bson:"voucher_id,omitempty"`
	Note          string             `bson:"note,omitempty"`
	CreatedAt     time.Time          `bson:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at"`
}

const (
	OrderStatusPending    = "pending"
	OrderStatusProcessing = "processing"
	OrderStatusShipping   = "shipping"
	OrderStatusDelivered  = "delivered"
	OrderStatusCanceled   = "canceled"
)

const (
	PaymentMethodCOD  = "cod"
	PaymentMethodBank = "bank"
)
