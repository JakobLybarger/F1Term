package models

type Event struct {
	Meeting   Meeting
	Session   Session
	Drivers   []Driver
	Positions []Position
}
