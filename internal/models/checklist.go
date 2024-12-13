package models

import (
	"context"
	"errors"
	"time"

	"github.com/timenglesf/bike-checkover-checklist/internal/validator"

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
	Description string              `bson:"description,omitempty"`
	Name        string              `bson:"-"`
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
	Notes         string        `bson:"notes,omitempty"`
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
		Notes:         "",
	}
}

type ChecklistDocument struct {
	Checklist Checklist          `bson:"checklist"`
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserId    primitive.ObjectID `bson:"userId"`
	StoreId   string             `bson:"storeId"`
	CreatedAt primitive.DateTime `bson:"createdAt"`
	Complete  bool               `bson:"complete"`
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
		Complete:  false,
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
	doc := createChecklistDocument(checklist, user.ID, user.StoreId)
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

func (m *ChecklistModel) updateChecklistDocument(
	ctx context.Context,
	documentId primitive.ObjectID,
	update bson.M,
) error {
	coll := m.getCollection()
	filter := bson.M{"_id": documentId}
	_, err := coll.UpdateOne(ctx, filter, update)
	return err
}

// Update updates the checklist property of a checklist document
// with the given documentId.
func (m *ChecklistModel) Update(
	ctx context.Context,
	documentId primitive.ObjectID,
	checklist Checklist,
) error {
	update := bson.M{
		"$set": bson.M{
			"checklist.brakePad.status":      checklist.BrakePad.Status,
			"checklist.chain.status":         checklist.Chain.Status,
			"checklist.tires.status":         checklist.Tires.Status,
			"checklist.cassette.status":      checklist.Cassette.Status,
			"checklist.cablesHousing.status": checklist.CablesHousing.Status,
			"checklist.tubes.status":         checklist.Tubes.Status,
			"checklist.chainRing.status":     checklist.ChainRing.Status,
			"checklist.frontWheel.status":    checklist.FrontWheel.Status,
			"checklist.padFunction.status":   checklist.PadFunction.Status,
			"checklist.derailleur.status":    checklist.Derailleur.Status,
			"checklist.rearWheel.status":     checklist.RearWheel.Status,
			"checklist.rotorRim.status":      checklist.RotorRim.Status,
			"checklist.hanger.status":        checklist.Hanger.Status,
			"checklist.shifting.status":      checklist.Shifting.Status,
			"checklist.notes":                checklist.Notes,
		},
	}
	return m.updateChecklistDocument(ctx, documentId, update)
}

// SubmitChecklist updates the checklist and
// sets complete fields of a checklist document to true.
func (m *ChecklistModel) SubmitChecklist(
	ctx context.Context,
	documentId primitive.ObjectID,
	checklist Checklist,
) error {
	update := bson.M{"$set": bson.M{"complete": true, "checklist": checklist}}
	return m.updateChecklistDocument(ctx, documentId, update)
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

func (m *ChecklistModel) GetRecentActiveChecklist(
	ctx context.Context,
	userId primitive.ObjectID,
) (*ChecklistDocument, error) {
	coll := m.getCollection()
	filter := bson.M{
		"userId":   userId,
		"complete": false,
	}
	var doc ChecklistDocument
	err := coll.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

// //////// Form //////////
type ChecklistForm struct {
	validator.Validator `form:"-"`
	BrakePad            string `form:"brake-pad"`
	Chain               string `form:"chain"`
	Tires               string `form:"tires"`
	Cassette            string `form:"cassette"`
	CablesHousing       string `form:"cables-housing"`
	Tubes               string `form:"tubes"`
	ChainRing           string `form:"chain-ring"`
	FrontWheel          string `form:"front-wheel"`
	PadFunction         string `form:"pad-function"`
	Derailleur          string `form:"derailleur"`
	RearWheel           string `form:"rear-wheel"`
	RotorRim            string `form:"rotor-rim"`
	Hanger              string `form:"hanger"`
	Shifting            string `form:"shifting"`
	Notes               string `form:"notes"`
}

func (cl ChecklistForm) ConvertFormToChecklist() Checklist {
	return Checklist{
		BrakePad:      ChecklistItem{Status: ChecklistItemStatus(cl.BrakePad), Name: "Brake Pad", Id: BrakePad},
		Chain:         ChecklistItem{Status: ChecklistItemStatus(cl.Chain), Name: "Chain", Id: Chain},
		Tires:         ChecklistItem{Status: ChecklistItemStatus(cl.Tires), Name: "Tires", Id: Tires},
		Cassette:      ChecklistItem{Status: ChecklistItemStatus(cl.Cassette), Name: "Cassette", Id: Cassette},
		CablesHousing: ChecklistItem{Status: ChecklistItemStatus(cl.CablesHousing), Name: "Cables Housing", Id: CablesHousing},
		Tubes:         ChecklistItem{Status: ChecklistItemStatus(cl.Tubes), Name: "Tubes", Id: Tubes},
		ChainRing:     ChecklistItem{Status: ChecklistItemStatus(cl.ChainRing), Name: "Chain Ring", Id: ChainRing},
		FrontWheel:    ChecklistItem{Status: ChecklistItemStatus(cl.FrontWheel), Name: "Front Wheel", Id: FrontWheel},
		PadFunction:   ChecklistItem{Status: ChecklistItemStatus(cl.PadFunction), Name: "Pad Function", Id: PadFunction},
		Derailleur:    ChecklistItem{Status: ChecklistItemStatus(cl.Derailleur), Name: "Derailleur", Id: Derailleur},
		RearWheel:     ChecklistItem{Status: ChecklistItemStatus(cl.RearWheel), Name: "Rear Wheel", Id: RearWheel},
		RotorRim:      ChecklistItem{Status: ChecklistItemStatus(cl.RotorRim), Name: "Rotor Rim", Id: RotorRim},
		Hanger:        ChecklistItem{Status: ChecklistItemStatus(cl.Hanger), Name: "Hanger", Id: Hanger},
		Shifting:      ChecklistItem{Status: ChecklistItemStatus(cl.Shifting), Name: "Shifting", Id: Shifting},
		Notes:         cl.Notes,
	}
}
