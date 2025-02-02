package logger

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Log struct {
	ID            primitive.ObjectID `bson:"_id"`
	UserID        uuid.UUID          `bson:"user_id,omitempty"`
	CorrelationID uuid.UUID          `bson:"correlation_id,omitempty"`
	RemoteIP      string             `bson:"remote_ip"`
	Referer       string             `bson:"referer"`
	UserAgent     string             `bson:"user_agent"`
	Method        string             `bson:"method"`
	Location      string             `bson:"location"`
	Message       string             `bson:"messages"`
	CreatedAt     time.Time          `bson:"created_at"`
	StatusCode    int                `bson:"status"`
	Headers       map[string]string  `bson:"headers"`
}
