package models

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Admin struct {
	FirstName string             `bson:"firstName"`
	LastName  string             `bson:"lastName"`
	Username  string             `bson:"username"`
	Password  string             `bson:"password"`
	CreatedAt primitive.DateTime `bson:"createdAt"`
	ID        primitive.ObjectID `bson:"_id,omitempty"`
}

type AdminModel struct {
	DB             *mongo.Client
	DBName         string
	CollectionName string
}

type AdminModelInterface interface {
	Insert(ctx context.Context, admin Admin) error
	GetDocumentCount(ctx context.Context) (int64, error)
	GetAdminByUsername(ctx context.Context, username string) (Admin, error)
}

func CreateAdmin(username, password, firstName, lastName string) (Admin, error) {
	hashedPassBS, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return Admin{}, err
	}
	return Admin{
		FirstName: firstName,
		LastName:  lastName,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		Username:  strings.ToLower(username),
		Password:  string(hashedPassBS),
	}, nil
}

func (m *AdminModel) getCollection() *mongo.Collection {
	return m.DB.Database(m.DBName).Collection(m.CollectionName)
}

func (m *AdminModel) Insert(
	ctx context.Context,
	admin Admin,
) error {
	coll := m.getCollection()
	_, err := coll.InsertOne(ctx, admin)
	if err != nil {
		return err
	}
	return nil
}

func (m *AdminModel) GetDocumentCount(ctx context.Context) (int64, error) {
	coll := m.getCollection()
	count, err := coll.CountDocuments(ctx, bson.D{})
	fmt.Println(count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (m *AdminModel) FindByUsername(ctx context.Context, username string) (Admin, error) {
	coll := m.getCollection()
	var admin Admin
	err := coll.FindOne(ctx, map[string]string{"username": strings.ToLower(username)}).Decode(&admin)
	if err != nil {
		return Admin{}, err
	}
	return admin, nil
}
