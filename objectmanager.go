package main

import (
	"errors"
	"log"
)

const (
	DragStatusIdle     = 0
	DragStatusDragging = 1
)

type CardObject struct {
	Position  [2]float64
	CID       string
	DragState DragState
}

type DragState struct {
	Status int
	UID    string
}

type ObjectManager struct {
	Cards        map[string]*CardObject
	UserDragging map[string]*CardObject
}

func NewObjectManager() *ObjectManager {
	om := &ObjectManager{
		Cards:        make(map[string]*CardObject),
		UserDragging: make(map[string]*CardObject),
	}
	om.Cards["aaa"] = &CardObject{
		Position: [2]float64{0., 0.},
		CID:      "aaa",
		DragState: DragState{
			DragStatusIdle,
			"",
		},
	}
	return om
}

func (om *ObjectManager) GetCard(CID string) (card *CardObject, err error) {
	card, ok := om.Cards[CID]
	if !ok {
		err = errors.New("not found")
		return
	}
	return
}

func (om *ObjectManager) StartDragging(UID string, CID string) error {
	card, err := om.GetCard(CID)
	if err != nil {
		return err
	}
	// check: is someone dragging the card
	if card.DragState.Status == DragStatusDragging {
		err = errors.New("someone is dragging the card")
		return err
	}
	// ready to drag
	card.DragState.Status = DragStatusDragging
	card.DragState.UID = UID
	return nil
}

func (om *ObjectManager) stopDragging(UID string, CID string) {
	delete(om.UserDragging, UID)
	card := om.Cards[CID]
	card.DragState.Status = DragStatusIdle
	card.DragState.UID = ""
}

func (om *ObjectManager) CancelDragging(UID string, CID string) error {
	card, err := om.GetCard(CID)
	if err != nil {
		return err
	}
	if card.DragState.Status != DragStatusDragging || card.DragState.UID != UID {
		return errors.New("no dragging to cancel")
	}
	om.stopDragging(UID, CID)
	return nil
}

func (om *ObjectManager) FinishDragging(UID string, CID string, pos [2]float64) error {
	card, err := om.GetCard(CID)
	if err != nil {
		return err
	}
	log.Println("dragging", card, UID)
	if card.DragState.Status != DragStatusDragging || card.DragState.UID != UID {
		return errors.New("no dragging to finish")
	}
	card.Position = pos
	om.stopDragging(UID, CID)
	return nil
}
