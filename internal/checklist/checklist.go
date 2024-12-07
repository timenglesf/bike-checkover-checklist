package checklist

type CheckListItemStatus string

type CheckListItemId string

var (
	BrakePad      CheckListItemId = "brake-pad"
	Chain         CheckListItemId = "chain"
	Tires         CheckListItemId = "tires"
	Cassette      CheckListItemId = "cassette"
	CablesHousing CheckListItemId = "cables-housing"
	Tubes         CheckListItemId = "tubes"
	ChainRing     CheckListItemId = "chain-ring"
	FrontWheel    CheckListItemId = "front-wheel"
	PadFunction   CheckListItemId = "pad-function"
	Derailleur    CheckListItemId = "derailleur"
	RearWheel     CheckListItemId = "rear-wheel"
	RotorRim      CheckListItemId = "rotor-rim"
	Hanger        CheckListItemId = "hanger"
	Shifting      CheckListItemId = "shifting"
)

type CheckListItem struct {
	Status      CheckListItemStatus
	Description string
	Name        string
	Id          CheckListItemId
}

var (
	Pass          CheckListItemStatus = "pass"
	Fail          CheckListItemStatus = "fail"
	NotApplicable CheckListItemStatus = "not-applicable"
)

func CreateCheckListItem(name string, id CheckListItemId, d string) CheckListItem {
	return CheckListItem{
		Status:      NotApplicable,
		Name:        name,
		Id:          id,
		Description: d,
	}
}

type CheckList struct {
	BrakePad      CheckListItem
	Chain         CheckListItem
	Tires         CheckListItem
	Cassette      CheckListItem
	CablesHousing CheckListItem
	Tubes         CheckListItem
	ChainRing     CheckListItem
	FrontWheel    CheckListItem
	PadFunction   CheckListItem
	Derailleur    CheckListItem
	RearWheel     CheckListItem
	RotorRim      CheckListItem
	Hanger        CheckListItem
	Shifting      CheckListItem
}

func CreateCheckList() *CheckList {
	return &CheckList{
		BrakePad:      CreateCheckListItem("Brake Pad", BrakePad, "Brake pad pass inspection"),
		Chain:         CreateCheckListItem("Chain", Chain, "Chain wear measurement pass inspection"),
		Tires:         CreateCheckListItem("Tires", Tires, "Tires pass inspection"),
		Cassette:      CreateCheckListItem("Cassette", Cassette, "Cassette pass inspection"),
		CablesHousing: CreateCheckListItem("Cables/Housing", CablesHousing, "Cables/Housing/Brake Hose pass inspection"),
		Tubes:         CreateCheckListItem("Tubes", Tubes, "Tubes pass inspection"),
		ChainRing:     CreateCheckListItem("Chain Ring", ChainRing, "Chairing wear measurement pass inspection"),
		FrontWheel:    CreateCheckListItem("Front Wheel", FrontWheel, "Front: spoke tension/wheel true pass inspection"),
		PadFunction:   CreateCheckListItem("Pad Function", PadFunction, "Test brake function and noise pass inspection"),
		Derailleur:    CreateCheckListItem("Derailleur", Derailleur, "Check derailluer for play and inspect pullies"),
		RearWheel:     CreateCheckListItem("Rear Wheel", RearWheel, "Rear: spoke tension/wheel true pass inspection"),
		RotorRim:      CreateCheckListItem("Rotor/Rim", RotorRim, "Rotor/Rim pass inspection"),
		Hanger:        CreateCheckListItem("Hanger", Hanger, "Hanger pass inspection"),
		Shifting:      CreateCheckListItem("Shifting", Shifting, "Shifting pass inspection"),
	}
}
