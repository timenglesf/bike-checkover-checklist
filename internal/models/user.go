package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/timenglesf/bike-checkover-checklist/internal/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const USER_COLL = "users"

type User struct {
	FirstName string             `bson:"firstName" form:"first-name"`
	LastName  string             `bson:"lastName" form:"last-name"`
	Pin       string             `bson:"pin" form:"pin"`
	StoreId   string             `bson:"storeId" form:"store-id"`
	CreatedAt primitive.DateTime `bson:"createdAt"`
	ID        primitive.ObjectID `bson:"_id,omitempty"`
}

type UserModel struct {
	DB             *mongo.Client
	DBName         string
	CollectionName string
}

type UserModelInterface interface {
	Insert(user User, ctx context.Context) error
	FindById(id primitive.ObjectID, ctx context.Context) (User, error)
	FindByPin(pin string, ctx context.Context) (User, error)
	FindAll(ctx context.Context) ([]User, error)
	FindUsersByStoreId(storeId string, ctx context.Context) ([]User, error)
	GetDocumentCountByPin(ctx context.Context, pin string) (int64, error)
}

func CreateUser(firstName, lastName, pin, storeId string) User {
	return User{
		FirstName: firstName,
		LastName:  lastName,
		Pin:       pin,
		StoreId:   storeId,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}
}

func (m *UserModel) getCollection() *mongo.Collection {
	return m.DB.Database(m.DBName).Collection(m.CollectionName)
}

func (m *UserModel) Insert(
	ctx context.Context,
	user User,
) error {
	coll := m.getCollection()
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

func (m *UserModel) FindById(
	ctx context.Context,
	id primitive.ObjectID,
) (User, error) {
	coll := m.getCollection()
	filter := bson.M{"_id": id}
	var user User
	err := coll.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (m *UserModel) FindByPin(
	ctx context.Context,
	pin string,
) (User, error) {
	fmt.Println("pin: ", pin)
	coll := m.getCollection()
	filter := bson.M{"pin": pin}
	var user User

	err := coll.FindOne(ctx, filter).
		Decode(&user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (m *UserModel) GetDocumentCountByPin(
	ctx context.Context,
	pin string,
) (int64, error) {
	coll := m.getCollection()
	filter := bson.M{"pin": pin}
	count, err := coll.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (m *UserModel) GetAll(ctx context.Context) ([]User, error) {
	coll := m.getCollection()
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

func (m *UserModel) GetUsersByStoreId(
	ctx context.Context,
	storeId string,
) ([]User, error) {
	coll := m.getCollection()
	filter := bson.M{"storeId": storeId}
	cursor, err := coll.Find(ctx, filter)
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

///////////////////////////////
/////// USER FORMS ////////////
///////////////////////////////

type UserForm struct {
	validator.Validator `form:"-"`
	User
	PinConfirm string `form:"pin-confirm"`
}
