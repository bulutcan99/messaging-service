package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// TODO: file attachment olcak mesaj ile file yukleyebilcek
// boylelikle cekerken file'i da cekebilcek
// TODO: hub room mantigi eklenecek ve http kullanilacak
// TODO: ws ile sadece notification ve aktiflik
// wse kullanici 15sn'de bir ping aticak
type MessageReceipt struct {
	UserID bson.ObjectID `bson:"user_id"`
	ReadAt time.Time     `bson:"read_at"`
}

type Message struct {
	ID        bson.ObjectID    `bson:"_id,omitempty"`
	SenderID  bson.ObjectID    `bson:"sender_id"`
	Content   string           `bson:"content"`
	ChatID    bson.ObjectID    `bson:"chat_id"`
	Receipts  []MessageReceipt `bson:"receipts"`
	CreatedAt time.Time        `bson:"created_at"`
	UpdatedAt time.Time        `bson:"updated_at"`
	DeletedAt *time.Time       `bson:"deleted_at,omitempty"`
}
