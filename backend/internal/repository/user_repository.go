package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"cinema-booking-backend/internal/model"
)

type UserRepository struct {
	users *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{users: db.Collection("users")}
}

func (r *UserRepository) UpsertGoogleUser(ctx context.Context, googleSub, email, name string) (*model.User, error) {
	now := time.Now()

	filter := bson.M{"google_sub": googleSub}
	update := bson.M{
		"$set": bson.M{
			"email": email,
			"name":  name,
		},
		"$setOnInsert": bson.M{
			"google_sub": googleSub,
			"role":       model.RoleUser,
			"created_at": now,
		},
	}
	opts := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)

	var user model.User
	err := r.users.FindOneAndUpdate(ctx, filter, update, opts).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}