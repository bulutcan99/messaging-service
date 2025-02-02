package entity

type Status int

const (
	Idle Status = iota
	Waiting
	Connected
	Disconnected
)

type TextEvent int

const (
	Typing TextEvent = iota
	Stopped
	Submitted
)

type ListenerEvent int

const (
	EventAdminAuth ListenerEvent = iota
	EventMessageError
	EventRoomRegister
	EventRoomUnregister
	EventRoomBroadcast
	EventRoomUpdate
	EventHubUpdate
)
