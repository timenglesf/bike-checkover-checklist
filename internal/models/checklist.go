package models

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChecklistItemStatus string

type ChecklistItemId string

var (
	BrakePad      ChecklistItemId = "brake-pad"
	Chain         ChecklistItemId = "chain"
	Tires         ChecklistItemId = "tires"
	Cassette      ChecklistItemId = "cassette"
	CablesHousing ChecklistItemId = "cables-housing"
	Tubes         ChecklistItemId = "tubes"
	ChainRing     ChecklistItemId = "chain-ring"
	FrontWheel    ChecklistItemId = "front-wheel"
	PadFunction   ChecklistItemId = "pad-function"
	Derailleur    ChecklistItemId = "derailleur"
	RearWheel     ChecklistItemId = "rear-wheel"
	RotorRim      ChecklistItemId = "rotor-rim"
	Hanger        ChecklistItemId = "hanger"
	Shifting      ChecklistItemId = "shifting"
)

type ChecklistItem struct {
	Status      ChecklistItemStatus `bson:"status"`
	Description string
	Name        string
	Id          ChecklistItemId
}

var (
	Pass          ChecklistItemStatus = "pass"
	Fail          ChecklistItemStatus = "fail"
	NotApplicable ChecklistItemStatus = "not-applicable"
)

func CreateChecklistItem(name string, id ChecklistItemId, d string) ChecklistItem {
	return ChecklistItem{
		Status:      NotApplicable,
		Name:        name,
		Id:          id,
		Description: d,
	}
}

type Checklist struct {
	BrakePad      ChecklistItem `bson:"brakePad"`
	Chain         ChecklistItem `bson:"chain"`
	Tires         ChecklistItem `bson:"tires"`
	Cassette      ChecklistItem `bson:"cassette"`
	CablesHousing ChecklistItem `bson:"cablesHousing"`
	Tubes         ChecklistItem `bson:"tubes"`
	ChainRing     ChecklistItem `bson:"chainRing"`
	FrontWheel    ChecklistItem `bson:"frontWheel"`
	PadFunction   ChecklistItem `bson:"padFunction"`
	Derailleur    ChecklistItem `bson:"derailleur"`
	RearWheel     ChecklistItem `bson:"rearWheel"`
	RotorRim      ChecklistItem `bson:"rotorRim"`
	Hanger        ChecklistItem `bson:"hanger"`
	Shifting      ChecklistItem `bson:"shifting"`
}

func CreateChecklist() *Checklist {
	return &Checklist{
		BrakePad:      CreateChecklistItem("Brake Pad", BrakePad, "Brake pad pass inspection"),
		Chain:         CreateChecklistItem("Chain", Chain, "Chain wear measurement pass inspection"),
		Tires:         CreateChecklistItem("Tires", Tires, "Tires pass inspection"),
		Cassette:      CreateChecklistItem("Cassette", Cassette, "Cassette pass inspection"),
		CablesHousing: CreateChecklistItem("Cables/Housing", CablesHousing, "Cables/Housing/Brake Hose pass inspection"),
		Tubes:         CreateChecklistItem("Tubes", Tubes, "Tubes pass inspection"),
		ChainRing:     CreateChecklistItem("Chain Ring", ChainRing, "Chairing wear measurement pass inspection"),
		FrontWheel:    CreateChecklistItem("Front Wheel", FrontWheel, "Front: spoke tension/wheel true pass inspection"),
		PadFunction:   CreateChecklistItem("Pad Function", PadFunction, "Test brake function and noise pass inspection"),
		Derailleur:    CreateChecklistItem("Derailleur", Derailleur, "Check derailluer for play and inspect pullies"),
		RearWheel:     CreateChecklistItem("Rear Wheel", RearWheel, "Rear: spoke tension/wheel true pass inspection"),
		RotorRim:      CreateChecklistItem("Rotor/Rim", RotorRim, "Rotor/Rim pass inspection"),
		Hanger:        CreateChecklistItem("Hanger", Hanger, "Hanger pass inspection"),
		Shifting:      CreateChecklistItem("Shifting", Shifting, "Shifting pass inspection"),
	}
}

type ChecklistDocument struct {
	Checklist Checklist          `bson:"checklist"`
	ID        primitive.ObjectID `bson:"_id"`
	UserId    primitive.ObjectID `bson:"userId"`
	StoreId   string             `bson:"storeId"`
	CreatedAt primitive.DateTime `bson:"createdAt"`
}

type ChecklistModel struct {
	DB             *mongo.Client
	DBName         string
	CollectionName string
}

// createChecklistDocument creates a ChecklistDocument from the given
// checklist, userId, and storeId.
func createChecklistDocument(
	checklist Checklist,
	userId primitive.ObjectID,
	storeId string,
) ChecklistDocument {
	return ChecklistDocument{
		Checklist: checklist,
		UserId:    userId,
		StoreId:   storeId,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}
}

// getCollection returns the MongoDB collection for checklists.
func (m *ChecklistModel) getCollection() *mongo.Collection {
	return m.DB.Database(m.DBName).Collection(m.CollectionName)
}

// Insert inserts a new checklist document into the database.
func (m *ChecklistModel) Insert(
	ctx context.Context,
	checklist Checklist,
	user User,
) error {
	coll := m.getCollection()
	doc := createChecklistDocument(
		checklist,
		user.ID,
		user.StoreId,
	)

	_, err := coll.InsertOne(ctx, doc)
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

// Update updates the checklist property of a checklist document
// with the given documentId.
func (m *ChecklistModel) Update(
	ctx context.Context,
	documentId primitive.ObjectID,
	checklist Checklist,
) error {
	coll := m.getCollection()
	filter := bson.M{"_id": documentId}
	update := bson.M{"$set": bson.M{"checklist": checklist}}
	_, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

// Get retrieves a checklist document by its documentId.
func (m *ChecklistModel) Get(
	ctx context.Context,
	documentId primitive.ObjectID,
) (*ChecklistDocument, error) {
	coll := m.getCollection()
	filter := bson.M{"_id": documentId}
	var doc ChecklistDocument
	err := coll.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

// GetUserChecklists retrieves all checklist documents for a given userId.
func (m *ChecklistModel) GetUserChecklists(
	ctx context.Context,
	userId primitive.ObjectID,
) ([]ChecklistDocument, error) {
	coll := m.getCollection()
	filter := bson.M{"userId": userId}
	opts := options.Find().SetSort(
		bson.D{{Key: "createdAt", Value: -1}},
	)
	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var checklists []ChecklistDocument
	err = cursor.All(ctx, &checklists)
	if err != nil {
		return nil, err
	}
	return checklists, nil
}

// GetStoreChecklists retrieves all checklist documents for a given storeId.
func (m *ChecklistModel) GetStoreChecklists(
	ctx context.Context,
	storeId string,
) ([]ChecklistDocument, error) {
	coll := m.getCollection()

	filter := bson.M{"storeId": storeId}
	opts := options.Find().SetSort(
		bson.D{{Key: "createdAt", Value: -1}},
	)
	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var checklists []ChecklistDocument
	err = cursor.All(ctx, &checklists)
	if err != nil {
		return nil, err
	}
	return checklists, nil
}
