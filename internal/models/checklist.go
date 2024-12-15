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

// ChecklistItemId and ChecklistItemStatus types
type (
	ChecklistItemId     string
	ChecklistItemStatus string
)

// ChecklistItemId constants
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

var (
	Pass          ChecklistItemStatus = "pass"
	Fail          ChecklistItemStatus = "fail"
	NotApplicable ChecklistItemStatus = "not-applicable"
)

// ChecklistItem struct
type ChecklistItem struct {
	Status      ChecklistItemStatus `bson:"status"`
	Description string              `bson:"description,omitempty"`
	Name        string              `bson:"-"`
	Id          ChecklistItemId
}

// Checklist struct
type Checklist struct {
	BrakePad      ChecklistItemStatus `bson:"brake-pad" form:"brake-pad"`
	Chain         ChecklistItemStatus `bson:"chain" form:"chain"`
	Tires         ChecklistItemStatus `bson:"tires" form:"tires"`
	Cassette      ChecklistItemStatus `bson:"cassette" form:"cassette"`
	CablesHousing ChecklistItemStatus `bson:"cables-housing" form:"cables-housing"`
	Tubes         ChecklistItemStatus `bson:"tubes" form:"tubes"`
	ChainRing     ChecklistItemStatus `bson:"chain-ring" form:"chain-ring"`
	FrontWheel    ChecklistItemStatus `bson:"front-wheel" form:"front-wheel"`
	PadFunction   ChecklistItemStatus `bson:"pad-function" form:"pad-function"`
	Derailleur    ChecklistItemStatus `bson:"derailleur" form:"derailleur"`
	RearWheel     ChecklistItemStatus `bson:"rear-wheel" form:"rear-wheel"`
	RotorRim      ChecklistItemStatus `bson:"rotor-rim" form:"rotor-rim"`
	Hanger        ChecklistItemStatus `bson:"hanger" form:"hanger"`
	Shifting      ChecklistItemStatus `bson:"shifting" form:"shifting"`
	Notes         string              `bson:"notes,omitempty" form:"notes"`
}

func CreateChecklist() *Checklist {
	return &Checklist{
		BrakePad:      NotApplicable,
		Chain:         NotApplicable,
		Tires:         NotApplicable,
		Cassette:      NotApplicable,
		CablesHousing: NotApplicable,
		Tubes:         NotApplicable,
		ChainRing:     NotApplicable,
		FrontWheel:    NotApplicable,
		PadFunction:   NotApplicable,
		Derailleur:    NotApplicable,
		RearWheel:     NotApplicable,
		RotorRim:      NotApplicable,
		Hanger:        NotApplicable,
		Shifting:      NotApplicable,
		Notes:         "",
	}
}

// BikeDescription struct
type BikeDescription struct {
	Brand string `bson:"brand" form:"brand"`
	Model string `bson:"model" form:"model"`
	Color string `bson:"color" form:"color"`
}

type ChecklistDocument struct {
	Checklist   Checklist          `bson:"checklist"`
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	UserId      primitive.ObjectID `bson:"userId"`
	StoreId     string             `bson:"storeId"`
	CreatedAt   primitive.DateTime `bson:"createdAt"`
	Complete    bool               `bson:"complete"`
	Description BikeDescription    `bson:"description,omitempty"`
}

// ChecklistModel struct
type ChecklistModel struct {
	DB             *mongo.Client
	DBName         string
	CollectionName string
}

// createChecklistDocument creates a ChecklistDocument from the given checklist, userId, and storeId.
func createChecklistDocument(checklist Checklist, userId primitive.ObjectID, storeId string) ChecklistDocument {
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
func (m *ChecklistModel) Insert(ctx context.Context, checklist Checklist, user User) error {
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

func (m *ChecklistModel) updateChecklistDocument(ctx context.Context, documentId primitive.ObjectID, update bson.M) error {
	coll := m.getCollection()
	filter := bson.M{"_id": documentId}
	_, err := coll.UpdateOne(ctx, filter, update)
	return err
}

func (m *ChecklistModel) Reset(ctx context.Context, documentId primitive.ObjectID) error {
	update := bson.M{"$set": bson.M{"complete": false, "checklist": CreateChecklist()}}
	return m.updateChecklistDocument(ctx, documentId, update)
}

// Update updates the checklist property of a checklist document with the given documentId.
func (m *ChecklistModel) Update(ctx context.Context, documentId primitive.ObjectID, checklist Checklist) error {
	update := bson.M{
		"$set": bson.M{
			"checklist.brake-pad":      checklist.BrakePad,
			"checklist.chain":          checklist.Chain,
			"checklist.tires":          checklist.Tires,
			"checklist.cassette":       checklist.Cassette,
			"checklist.cables-housing": checklist.CablesHousing,
			"checklist.tubes":          checklist.Tubes,
			"checklist.chain-ring":     checklist.ChainRing,
			"checklist.front-wheel":    checklist.FrontWheel,
			"checklist.pad-function":   checklist.PadFunction,
			"checklist.derailleur":     checklist.Derailleur,
			"checklist.rear-wheel":     checklist.RearWheel,
			"checklist.rotor-rim":      checklist.RotorRim,
			"checklist.hanger":         checklist.Hanger,
			"checklist.shifting":       checklist.Shifting,
			"checklist.notes":          checklist.Notes,
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
	description BikeDescription,
) error {
	update := bson.M{"$set": bson.M{"complete": true, "checklist": checklist, "description": description}}
	return m.updateChecklistDocument(ctx, documentId, update)
}

// Get retrieves a checklist document by its documentId.
func (m *ChecklistModel) Get(ctx context.Context, documentId primitive.ObjectID) (*ChecklistDocument, error) {
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
func (m *ChecklistModel) GetUserChecklists(ctx context.Context, userId primitive.ObjectID) ([]ChecklistDocument, error) {
	coll := m.getCollection()
	filter := bson.M{"userId": userId, "complete": true}
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
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
func (m *ChecklistModel) GetStoreChecklists(ctx context.Context, storeId string) ([]ChecklistDocument, error) {
	coll := m.getCollection()
	filter := bson.M{"storeId": storeId}
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
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

func (m *ChecklistModel) GetRecentActiveChecklist(ctx context.Context, userId primitive.ObjectID) (*ChecklistDocument, error) {
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

// ChecklistDisplay struct
type ChecklistDisplay struct {
	BrakePad      ChecklistItem
	Chain         ChecklistItem
	Tires         ChecklistItem
	Cassette      ChecklistItem
	CablesHousing ChecklistItem
	Tubes         ChecklistItem
	ChainRing     ChecklistItem
	FrontWheel    ChecklistItem
	PadFunction   ChecklistItem
	Derailleur    ChecklistItem
	RearWheel     ChecklistItem
	RotorRim      ChecklistItem
	Hanger        ChecklistItem
	Shifting      ChecklistItem
	Notes         string
	BikeDescription
}

func CreateChecklistDisplayItem(name string, id ChecklistItemId, d string) ChecklistItem {
	return ChecklistItem{
		Status:      NotApplicable,
		Name:        name,
		Id:          id,
		Description: d,
	}
}

func CreateChecklistDisplay() *ChecklistDisplay {
	return &ChecklistDisplay{
		BrakePad:        CreateChecklistDisplayItem("Brake Pad", BrakePad, "Brake pad pass inspection"),
		Chain:           CreateChecklistDisplayItem("Chain", Chain, "Chain wear measurement pass inspection"),
		Tires:           CreateChecklistDisplayItem("Tires", Tires, "Tires pass inspection"),
		Cassette:        CreateChecklistDisplayItem("Cassette", Cassette, "Cassette pass inspection"),
		CablesHousing:   CreateChecklistDisplayItem("Cables/Housing", CablesHousing, "Cables/Housing/Brake Hose pass inspection"),
		Tubes:           CreateChecklistDisplayItem("Tubes", Tubes, "Tubes pass inspection"),
		ChainRing:       CreateChecklistDisplayItem("Chain Ring", ChainRing, "Chairing wear measurement pass inspection"),
		FrontWheel:      CreateChecklistDisplayItem("Front Wheel", FrontWheel, "Front: spoke tension/wheel true pass inspection"),
		PadFunction:     CreateChecklistDisplayItem("Pad Function", PadFunction, "Test brake function and noise pass inspection"),
		Derailleur:      CreateChecklistDisplayItem("Derailleur", Derailleur, "Check derailluer for play and inspect pullies"),
		RearWheel:       CreateChecklistDisplayItem("Rear Wheel", RearWheel, "Rear: spoke tension/wheel true pass inspection"),
		RotorRim:        CreateChecklistDisplayItem("Rotor/Rim", RotorRim, "Rotor/Rim pass inspection"),
		Hanger:          CreateChecklistDisplayItem("Hanger", Hanger, "Hanger pass inspection"),
		Shifting:        CreateChecklistDisplayItem("Shifting", Shifting, "Shifting pass inspection"),
		Notes:           "",
		BikeDescription: BikeDescription{},
	}
}

func (clDisplay *ChecklistDisplay) ExtractChecklist() *Checklist {
	return &Checklist{
		BrakePad:      clDisplay.BrakePad.Status,
		Chain:         clDisplay.Chain.Status,
		Tires:         clDisplay.Tires.Status,
		Cassette:      clDisplay.Cassette.Status,
		CablesHousing: clDisplay.CablesHousing.Status,
		Tubes:         clDisplay.Tubes.Status,
		ChainRing:     clDisplay.ChainRing.Status,
		FrontWheel:    clDisplay.FrontWheel.Status,
		PadFunction:   clDisplay.PadFunction.Status,
		Derailleur:    clDisplay.Derailleur.Status,
		RearWheel:     clDisplay.RearWheel.Status,
		RotorRim:      clDisplay.RotorRim.Status,
		Hanger:        clDisplay.Hanger.Status,
		Shifting:      clDisplay.Shifting.Status,
		Notes:         "",
	}
}

func (clDisplay *ChecklistDisplay) UpdateStatusFromChecklist(checklist Checklist) {
	clDisplay.BrakePad.Status = checklist.BrakePad
	clDisplay.Chain.Status = checklist.Chain
	clDisplay.Tires.Status = checklist.Tires
	clDisplay.Cassette.Status = checklist.Cassette
	clDisplay.CablesHousing.Status = checklist.CablesHousing
	clDisplay.Tubes.Status = checklist.Tubes
	clDisplay.ChainRing.Status = checklist.ChainRing
	clDisplay.FrontWheel.Status = checklist.FrontWheel
	clDisplay.PadFunction.Status = checklist.PadFunction
	clDisplay.Derailleur.Status = checklist.Derailleur
	clDisplay.RearWheel.Status = checklist.RearWheel
	clDisplay.RotorRim.Status = checklist.RotorRim
	clDisplay.Hanger.Status = checklist.Hanger
	clDisplay.Shifting.Status = checklist.Shifting
	clDisplay.Notes = checklist.Notes
}

// ChecklistForm struct
type ChecklistForm struct {
	validator.Validator `form:"-"`
	Checklist
	BikeDescription
}

func (cl ChecklistForm) ConvertFormToChecklist() Checklist {
	return Checklist{
		BrakePad:      ChecklistItemStatus(cl.BrakePad),
		Chain:         ChecklistItemStatus(cl.Chain),
		Tires:         ChecklistItemStatus(cl.Tires),
		Cassette:      ChecklistItemStatus(cl.Cassette),
		CablesHousing: ChecklistItemStatus(cl.CablesHousing),
		Tubes:         ChecklistItemStatus(cl.Tubes),
		ChainRing:     ChecklistItemStatus(cl.ChainRing),
		FrontWheel:    ChecklistItemStatus(cl.FrontWheel),
		PadFunction:   ChecklistItemStatus(cl.PadFunction),
		Derailleur:    ChecklistItemStatus(cl.Derailleur),
		RearWheel:     ChecklistItemStatus(cl.RearWheel),
		RotorRim:      ChecklistItemStatus(cl.RotorRim),
		Hanger:        ChecklistItemStatus(cl.Hanger),
		Shifting:      ChecklistItemStatus(cl.Shifting),
		Notes:         cl.Notes,
	}
}

func (cl ChecklistForm) ConvertFormToBikeDescription() BikeDescription {
	return BikeDescription{
		Brand: cl.Brand,
		Model: cl.Model,
		Color: cl.Color,
	}
}
