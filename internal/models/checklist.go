package models

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
}

type ChecklistModel struct {
	DB             *mongo.Client
	DBName         string
	CollectionName string
}

func (m *ChecklistModel) Insert(
	ctx context.Context,
	checkList Checklist,
	userId primitive.ObjectID,
) error {
	coll := m.DB.Database(m.DBName).Collection(m.CollectionName)
	_, err := coll.InsertOne(ctx, checkList)
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
