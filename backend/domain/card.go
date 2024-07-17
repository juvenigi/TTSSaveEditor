package domain

import (
	"fmt"
	"strconv"
)

type Transform struct {
	PosX   float64 `json:"posX"`
	PosY   float64 `json:"posY"`
	PosZ   float64 `json:"posZ"`
	RotX   float64 `json:"rotX"`
	RotY   float64 `json:"rotY"`
	RotZ   float64 `json:"rotZ"`
	ScaleX float64 `json:"scaleX"`
	ScaleY float64 `json:"scaleY"`
	ScaleZ float64 `json:"scaleZ"`
}

type AltLookAngle struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type ColorDiffuse struct {
	R float64 `json:"r"`
	G float64 `json:"g"`
	B float64 `json:"b"`
}

type CustomDeck struct {
	FaceURL      string `json:"FaceURL"`
	BackURL      string `json:"BackURL"`
	NumWidth     int    `json:"NumWidth"`
	NumHeight    int    `json:"NumHeight"`
	BackIsHidden bool   `json:"BackIsHidden"`
	UniqueBack   bool   `json:"UniqueBack"`
	Type         int    `json:"Type"`
}

type Card struct {
	GUID                 string                `json:"GUID"`
	Name                 string                `json:"Name"`
	Transform            Transform             `json:"Transform"`
	Nickname             string                `json:"Nickname"`
	Description          string                `json:"Description"`
	GMNotes              string                `json:"GMNotes"`
	AltLookAngle         AltLookAngle          `json:"AltLookAngle"`
	ColorDiffuse         ColorDiffuse          `json:"ColorDiffuse"`
	LayoutGroupSortIndex int                   `json:"LayoutGroupSortIndex"`
	Value                int                   `json:"Value"`
	Locked               bool                  `json:"Locked"`
	Grid                 bool                  `json:"Grid"`
	Snap                 bool                  `json:"Snap"`
	IgnoreFoW            bool                  `json:"IgnoreFoW"`
	MeasureMovement      bool                  `json:"MeasureMovement"`
	DragSelectable       bool                  `json:"DragSelectable"`
	Autoraise            bool                  `json:"Autoraise"`
	Sticky               bool                  `json:"Sticky"`
	Tooltip              bool                  `json:"Tooltip"`
	GridProjection       bool                  `json:"GridProjection"`
	HideWhenFaceDown     bool                  `json:"HideWhenFaceDown"`
	Hands                bool                  `json:"Hands"`
	CardID               int                   `json:"CardID"`
	SidewaysCard         bool                  `json:"SidewaysCard"`
	CustomDeck           map[string]CustomDeck `json:"CustomDeck"`
	LuaScript            string                `json:"LuaScript"`
	LuaScriptState       string                `json:"LuaScriptState"`
	XmlUI                string                `json:"XmlUI"`
}

type PartialCard struct {
	LuaScript      string `json:"LuaScript"`
	LuaScriptState string `json:"LuaScriptState"`
}

type Deck struct {
	GUID                 string                `json:"GUID"`
	Name                 string                `json:"Name"`
	Transform            Transform             `json:"Transform"`
	Nickname             string                `json:"Nickname"`
	Description          string                `json:"Description"`
	GMNotes              string                `json:"GMNotes"`
	AltLookAngle         AltLookAngle          `json:"AltLookAngle"`
	ColorDiffuse         ColorDiffuse          `json:"ColorDiffuse"`
	LayoutGroupSortIndex int                   `json:"LayoutGroupSortIndex"`
	Value                int                   `json:"Value"`
	Locked               bool                  `json:"Locked"`
	Grid                 bool                  `json:"Grid"`
	Snap                 bool                  `json:"Snap"`
	IgnoreFoW            bool                  `json:"IgnoreFoW"`
	MeasureMovement      bool                  `json:"MeasureMovement"`
	DragSelectable       bool                  `json:"DragSelectable"`
	Autoraise            bool                  `json:"Autoraise"`
	Sticky               bool                  `json:"Sticky"`
	Tooltip              bool                  `json:"Tooltip"`
	GridProjection       bool                  `json:"GridProjection"`
	HideWhenFaceDown     bool                  `json:"HideWhenFaceDown"`
	Hands                bool                  `json:"Hands"`
	SidewaysCard         bool                  `json:"SidewaysCard"`
	DeckIDs              []int                 `json:"DeckIDs"`
	CustomDeck           map[string]CustomDeck `json:"CustomDeck"`
	LuaScript            string                `json:"LuaScript"`
	LuaScriptState       string                `json:"LuaScriptState"`
	XmlUI                string                `json:"XmlUI"`
	ContainedObjects     []Card                `json:"ContainedObjects"` // Assuming the objects are of type Card
}

func (deckObject *Deck) AppendNewCard(luaScript string, luaScriptState string) error {
	var err error
	var maxCardID = len(deckObject.DeckIDs)

	// obtain DeckID prefix
	var deckKey string
	for key := range deckObject.CustomDeck {
		deckKey = key
		break
	}
	// effectively clone the first newCard in the deck to obtain the reference newCard
	var referenceCard Card
	err, _ = remarshal(deckObject.ContainedObjects[0], &referenceCard)
	if err != nil {
		return err
	}
	// now update GUID, CardID, and DeckID of the newCard
	referenceCard.GUID = "000000" //TODO
	cardIdNum, _ := strconv.Atoi(fmt.Sprintf("%s%02d", deckKey, maxCardID))
	referenceCard.CardID = cardIdNum
	// the values we care about
	referenceCard.LuaScript = luaScript
	referenceCard.LuaScriptState = luaScriptState
	// housekeeping
	deckObject.DeckIDs = append(deckObject.DeckIDs, cardIdNum)
	deckObject.ContainedObjects = append(deckObject.ContainedObjects, referenceCard)
	return nil
}
