package repository

import (
	"context"
	"messaging-service/internal/models"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ChatRepository struct {
	collection *mongo.Collection
}

func NewChatRepository(collection *mongo.Collection) *ChatRepository {
	return &ChatRepository{
		collection: collection,
	}
}

func (r *ChatRepository) GetChatByID(ctx context.Context, id string) (*models.Chat, error) {
	var chat models.Chat
	err := r.collection.FindOne(ctx, map[string]string{"chat_id": id}).Decode(&chat)
	if err != nil {
		return nil, err
	}
	return &chat, nil
}

func (r *ChatRepository) GetChatsByUserID(ctx context.Context, userID string) ([]models.Chat, error) {
	var chats []models.Chat
	cursor, err := r.collection.Find(ctx, map[string]string{"user_id": userID})
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &chats)
	if err != nil {
		return nil, err
	}
	return chats, nil
}
