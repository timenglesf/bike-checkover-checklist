package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const USER_COLL = "users"

type User struct {
	FirstName string             `bson:"firstName"`
	LastName  string             `bson:"lastName"`
	CreatedAt time.Time          `bson:"created_at"`
	ID        primitive.ObjectID `bson:"_id"`
	Pin       string             `bson:"pin"`
	StoreId   string             `bson:"storeId"`
}

type UserModel struct {
	DB             *mongo.Client
	DBName         string
	CollectionName string
}

func (m *UserModel) Insert(user User, ctx context.Context) error {
	coll := m.DB.Database(m.DBName).Collection(m.CollectionName)
	_, err := coll.InsertOne(ctx, user)
	if err != nil {
		var writeException mongo.WriteException
		if errors.As(err, &writeException) {
			for _, writeError := range writeException.WriteErrors {
				if writeError.Code == 11000 {
					return ErrDuplicate
				}
			}
		}
		return err
	}
	return nil
}

func (m *UserModel) FindByID(id primitive.ObjectID, ctx context.Context) (User, error) {
	coll := m.DB.Database(m.DBName).Collection(m.CollectionName)
	var user User
	err := coll.FindOne(ctx, User{ID: id}).Decode(&user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (m *UserModel) FindByPin(pin string, ctx context.Context) (User, error) {
	fmt.Println("pin: ", pin)
	coll := m.DB.Database(m.DBName).Collection(m.CollectionName)
	filter := bson.M{"pin": pin}
	var user User

	err := coll.FindOne(ctx, filter).
		Decode(&user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (m *UserModel) GetAll(ctx context.Context) ([]User, error) {
	coll := m.DB.Database(m.DBName).Collection(m.CollectionName)
	cursor, err := coll.Find(ctx, nil)
	if err != nil {
		return nil, err
	}
	var users []User
	err = cursor.All(ctx, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}
