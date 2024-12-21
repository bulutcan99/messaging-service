package chat

type Status int

const (
	Waiting Status = iota
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
