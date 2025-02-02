// Go
package http

import (
	"net/http"
	"time"
	"websocket-azure/internal/core/entity"
	"websocket-azure/internal/presentation/http/controller/_default"
	"websocket-azure/internal/presentation/http/controller/websocket"

	"github.com/google/uuid"
)

func (s *Server) SetupRouter() {
	s.websocket()
	s.index()
	s.message()
	s.conversations()
	s.users()
	s.messages()
}
func (s *Server) websocket() {
	// GET /ws
	s.app.HandleFunc("/ws", websocket.WsHandler).Methods(http.MethodGet)
}
func (s *Server) index() {
	s.app.HandleFunc("/health", _default.HealthCheckHandler).Methods(http.MethodGet)
}

func (s *Server) message() {
	s.app.HandleFunc("/conversations/{conversationId}/messages", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Hello, World!"}`))
		//I need to send this message to users connected to my websocket here
	}).Methods(http.MethodPost)
}

func (s *Server) conversations() {
	// GET /conversations/{conversationId}
	s.app.HandleFunc("/conversations/{conversationId}", func(w http.ResponseWriter, r *http.Request) {
		// TODO: write handler logic
	}).Methods(http.MethodGet)

	// GET /conversations/search
	s.app.HandleFunc("/conversations/search", func(w http.ResponseWriter, r *http.Request) {
		// TODO: write handler logic
	}).Methods(http.MethodGet)

	// GET /conversations/{conversationId}/search
	s.app.HandleFunc("/conversations/{conversationId}/search", func(w http.ResponseWriter, r *http.Request) {
		// TODO: write handler logic
	}).Methods(http.MethodGet)
}

func (s *Server) users() {
	// GET /users/{userId}/conversations
	s.app.HandleFunc("/users/{userId}/conversations", func(w http.ResponseWriter, r *http.Request) {
		// TODO: write handler logic
	}).Methods(http.MethodGet)
}

func (s *Server) messages() {
	// POST /messages/{userId}/read
	s.app.HandleFunc("/messages/{userId}/read", func(w http.ResponseWriter, r *http.Request) {
		// TODO: write handler logic
	}).Methods(http.MethodPost)

	// POST /messages/{userId}/delivered
	s.app.HandleFunc("/messages/{userId}/delivered", func(w http.ResponseWriter, r *http.Request) {
		// TODO: write handler logic
	}).Methods(http.MethodPost)
}

// --------------------------------------------------------------------
// Users
// --------------------------------------------------------------------
var MockUsers = []entity.User{
	{
		ID:                    uuid.MustParse("a7d3b2c4-3e6e-4b57-9c8d-3d9d1efae3c5"),
		BoardID:               "board-1",
		Username:              "john_doe",
		Email:                 "john@example.com",
		IsEmailVerified:       true,
		SchoolEmail:           "john.doe@school.com",
		IsSchoolEmailVerified: false,
		PhoneNumber:           "+1234567890",
		IsPhoneNumberVerified: true,
		Password:              "hashed_password",
		Avatar:                "https://example.com/avatars/john.png",
		Role:                  "Student",
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
		DeletedAt:             nil,
	},
	{
		ID:                    uuid.MustParse("b8c4d3e5-2f7f-4d8a-9b6a-5f7d1e2a3c4b"),
		BoardID:               "board-2",
		Username:              "jane_doe",
		Email:                 "jane@example.com",
		IsEmailVerified:       true,
		SchoolEmail:           "jane.doe@school.com",
		IsSchoolEmailVerified: true,
		PhoneNumber:           "+1987654321",
		IsPhoneNumberVerified: true,
		Password:              "hashed_password",
		Avatar:                "https://example.com/avatars/jane.png",
		Role:                  "Student",
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
		DeletedAt:             nil,
	},
	{
		ID:                    uuid.MustParse("c9d5e6f7-3a8b-4c9d-8e7f-6a5b4c3d2e1f"),
		BoardID:               "board-3",
		Username:              "alice_smith",
		Email:                 "alice@example.com",
		IsEmailVerified:       false,
		SchoolEmail:           "alice.smith@school.com",
		IsSchoolEmailVerified: false,
		PhoneNumber:           "+1122334455",
		IsPhoneNumberVerified: false,
		Password:              "hashed_password",
		Avatar:                "https://example.com/avatars/alice.png",
		Role:                  "Student",
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
		DeletedAt:             nil,
	},
	{
		ID:                    uuid.MustParse("d1e2f3a4-5b6c-7d8e-9f0a-1b2c3d4e5f6a"),
		BoardID:               "board-4",
		Username:              "bob_jones",
		Email:                 "bob@example.com",
		IsEmailVerified:       true,
		SchoolEmail:           "bob.jones@school.com",
		IsSchoolEmailVerified: true,
		PhoneNumber:           "+1223344556",
		IsPhoneNumberVerified: false,
		Password:              "hashed_password",
		Avatar:                "https://example.com/avatars/bob.png",
		Role:                  "Student",
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
		DeletedAt:             nil,
	},
	{
		ID:                    uuid.MustParse("e2f3a4b5-6c7d-8e9f-0a1b-2c3d4e5f6a7b"),
		BoardID:               "board-5",
		Username:              "charlie_kelly",
		Email:                 "charlie@example.com",
		IsEmailVerified:       false,
		SchoolEmail:           "charlie.kelly@school.com",
		IsSchoolEmailVerified: false,
		PhoneNumber:           "+1555666777",
		IsPhoneNumberVerified: true,
		Password:              "hashed_password",
		Avatar:                "https://example.com/avatars/charlie.png",
		Role:                  "Student",
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
		DeletedAt:             nil,
	},
}

// --------------------------------------------------------------------
// Conversations
// (For clarity we use “deterministic” UUIDs here so that each direct
// conversation gets its own unique ID rather than re‐using a user ID.)
// --------------------------------------------------------------------
var MockConversations = []entity.Conversation{
	{
		ID:               uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		ConversationType: "direct",
		Participants: []string{
			"a7d3b2c4-3e6e-4b57-9c8d-3d9d1efae3c5", // john_doe
			"b8c4d3e5-2f7f-4d8a-9b6a-5f7d1e2a3c4b", // jane_doe
		},
		GroupName: "",
	},
	{
		ID:               uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		ConversationType: "direct",
		Participants: []string{
			"b8c4d3e5-2f7f-4d8a-9b6a-5f7d1e2a3c4b", // jane_doe
			"c9d5e6f7-3a8b-4c9d-8e7f-6a5b4c3d2e1f", // alice_smith
		},
		GroupName: "",
	},
	{
		ID:               uuid.MustParse("33333333-3333-3333-3333-333333333333"),
		ConversationType: "direct",
		Participants: []string{
			"c9d5e6f7-3a8b-4c9d-8e7f-6a5b4c3d2e1f", // alice_smith
			"d1e2f3a4-5b6c-7d8e-9f0a-1b2c3d4e5f6a", // bob_jones
		},
		GroupName: "",
	},
	{
		ID:               uuid.MustParse("44444444-4444-4444-4444-444444444444"),
		ConversationType: "direct",
		Participants: []string{
			"d1e2f3a4-5b6c-7d8e-9f0a-1b2c3d4e5f6a", // bob_jones
			"e2f3a4b5-6c7d-8e9f-0a1b-2c3d4e5f6a7b", // charlie_kelly
		},
		GroupName: "",
	},
	{
		ID:               uuid.MustParse("55555555-5555-5555-5555-555555555555"),
		ConversationType: "direct",
		Participants: []string{
			"e2f3a4b5-6c7d-8e9f-0a1b-2c3d4e5f6a7b", // charlie_kelly
			"a7d3b2c4-3e6e-4b57-9c8d-3d9d1efae3c5", // john_doe
		},
		GroupName: "",
	},
	{
		ID:               uuid.MustParse("66666666-6666-6666-6666-666666666666"),
		ConversationType: "group",
		Participants: []string{
			"a7d3b2c4-3e6e-4b57-9c8d-3d9d1efae3c5", // john_doe
			"b8c4d3e5-2f7f-4d8a-9b6a-5f7d1e2a3c4b", // jane_doe
			"c9d5e6f7-3a8b-4c9d-8e7f-6a5b4c3d2e1f", // alice_smith
		},
		Admins:    []string{"a7d3b2c4-3e6e-4b57-9c8d-3d9d1efae3c5"},
		GroupName: "Fast and Furious",
	},
	{
		ID:               uuid.MustParse("77777777-7777-7777-7777-777777777777"),
		ConversationType: "group",
		Participants: []string{
			"d1e2f3a4-5b6c-7d8e-9f0a-1b2c3d4e5f6a", // bob_jones
			"e2f3a4b5-6c7d-8e9f-0a1b-2c3d4e5f6a7b", // charlie_kelly
			"a7d3b2c4-3e6e-4b57-9c8d-3d9d1efae3c5", // john_doe
		},
		Admins:    []string{"d1e2f3a4-5b6c-7d8e-9f0a-1b2c3d4e5f6a"},
		GroupName: "Matrix",
	},
	{
		ID:               uuid.MustParse("88888888-8888-8888-8888-888888888888"),
		ConversationType: "group",
		Participants: []string{
			"e2f3a4b5-6c7d-8e9f-0a1b-2c3d4e5f6a7b", // charlie_kelly
			"a7d3b2c4-3e6e-4b57-9c8d-3d9d1efae3c5", // john_doe
			"b8c4d3e5-2f7f-4d8a-9b6a-5f7d1e2a3c4b", // jane_doe
		},
		Admins:    []string{"e2f3a4b5-6c7d-8e9f-0a1b-2c3d4e5f6a7b", "a7d3b2c4-3e6e-4b57-9c8d-3d9d1efae3c5"},
		GroupName: "Lord of the Rings",
	},
}

// --------------------------------------------------------------------
// Messages (Here we create eight messages in the direct conversation
// between john_doe and jane_doe – i.e. conversation ID 11111111-1111-1111-1111-111111111111)
// --------------------------------------------------------------------
var MockMessages = []entity.Message{
	{
		ID:             uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
		ConversationID: uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		Message:        "Hello, how are you?",
		SenderUserID:   "a7d3b2c4-3e6e-4b57-9c8d-3d9d1efae3c5",
		File:           "",
		FileName:       "",
		FileType:       "",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	},
	{
		ID:             uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"),
		ConversationID: uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		Message:        "I'm good, thank you!",
		SenderUserID:   "b8c4d3e5-2f7f-4d8a-9b6a-5f7d1e2a3c4b",
		File:           "",
		FileName:       "",
		FileType:       "",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	},
	{
		ID:             uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc"),
		ConversationID: uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		Message:        "I'm not feeling well today.",
		SenderUserID:   "a7d3b2c4-3e6e-4b57-9c8d-3d9d1efae3c5",
		File:           "",
		FileName:       "",
		FileType:       "",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	},
	{
		ID:             uuid.MustParse("dddddddd-dddd-dddd-dddd-dddddddddddd"),
		ConversationID: uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		Message:        "I'm sorry to hear that. Get well soon!",
		SenderUserID:   "b8c4d3e5-2f7f-4d8a-9b6a-5f7d1e2a3c4b",
		File:           "",
		FileName:       "",
		FileType:       "",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	},
	{
		ID:             uuid.MustParse("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"),
		ConversationID: uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		Message:        "I'm feeling better now. Thanks!",
		SenderUserID:   "a7d3b2c4-3e6e-4b57-9c8d-3d9d1efae3c5",
		File:           "",
		FileName:       "",
		FileType:       "",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	},
	{
		ID:             uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff"),
		ConversationID: uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		Message:        "I'm glad to hear that. Let's meet up for lunch.",
		SenderUserID:   "b8c4d3e5-2f7f-4d8a-9b6a-5f7d1e2a3c4b",
		File:           "",
		FileName:       "",
		FileType:       "",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	},
	{
		ID:             uuid.MustParse("11111111-2222-3333-4444-555555555555"),
		ConversationID: uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		Message:        "Sure, let's meet at 12:30 PM.",
		SenderUserID:   "a7d3b2c4-3e6e-4b57-9c8d-3d9d1efae3c5",
		File:           "",
		FileName:       "",
		FileType:       "",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	},
	{
		ID:             uuid.MustParse("66666666-7777-8888-9999-000000000000"),
		ConversationID: uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		Message:        "Great! I'll see you then.",
		SenderUserID:   "b8c4d3e5-2f7f-4d8a-9b6a-5f7d1e2a3c4b",
		File:           "",
		FileName:       "",
		FileType:       "",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	},
}

// --------------------------------------------------------------------
// Conversation (Message) Events for the messages above.
// We make sure to include only the other participant from the conversation.
// (In this example the direct conversation “1111…” is between john_doe and jane_doe.)
// --------------------------------------------------------------------
var MockConversationEvents = []entity.MessageEvent{
	{
		ID:        uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
		MessageID: uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
		Delivered: []entity.Event{
			{
				ParticipantID: uuid.MustParse("b8c4d3e5-2f7f-4d8a-9b6a-5f7d1e2a3c4b"),
				Timestamp:     time.Now(),
			},
		},
		Read:      []entity.Event{}, // not yet read by jane_doe
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:        uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"),
		MessageID: uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"),
		Delivered: []entity.Event{
			{
				ParticipantID: uuid.MustParse("a7d3b2c4-3e6e-4b57-9c8d-3d9d1efae3c5"),
				Timestamp:     time.Now(),
			},
		},
		Read: []entity.Event{
			{
				ParticipantID: uuid.MustParse("a7d3b2c4-3e6e-4b57-9c8d-3d9d1efae3c5"),
				Timestamp:     time.Now(),
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:        uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc"),
		MessageID: uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc"),
		Delivered: []entity.Event{
			{
				ParticipantID: uuid.MustParse("b8c4d3e5-2f7f-4d8a-9b6a-5f7d1e2a3c4b"),
				Timestamp:     time.Now(),
			},
		},
		Read:      []entity.Event{}, // not yet read
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:        uuid.MustParse("dddddddd-dddd-dddd-dddd-dddddddddddd"),
		MessageID: uuid.MustParse("dddddddd-dddd-dddd-dddd-dddddddddddd"),
		Delivered: []entity.Event{
			{
				ParticipantID: uuid.MustParse("a7d3b2c4-3e6e-4b57-9c8d-3d9d1efae3c5"),
				Timestamp:     time.Now(),
			},
		},
		Read: []entity.Event{
			{
				ParticipantID: uuid.MustParse("a7d3b2c4-3e6e-4b57-9c8d-3d9d1efae3c5"),
				Timestamp:     time.Now(),
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:        uuid.MustParse("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"),
		MessageID: uuid.MustParse("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"),
		Delivered: []entity.Event{
			{
				ParticipantID: uuid.MustParse("b8c4d3e5-2f7f-4d8a-9b6a-5f7d1e2a3c4b"),
				Timestamp:     time.Now(),
			},
		},
		Read: []entity.Event{
			{
				ParticipantID: uuid.MustParse("b8c4d3e5-2f7f-4d8a-9b6a-5f7d1e2a3c4b"),
				Timestamp:     time.Now(),
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:        uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff"),
		MessageID: uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff"),
		Delivered: []entity.Event{
			{
				ParticipantID: uuid.MustParse("a7d3b2c4-3e6e-4b57-9c8d-3d9d1efae3c5"),
				Timestamp:     time.Now(),
			},
		},
		Read:      []entity.Event{}, // not yet read
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:        uuid.MustParse("11111111-2222-3333-4444-555555555555"),
		MessageID: uuid.MustParse("11111111-2222-3333-4444-555555555555"),
		Delivered: []entity.Event{
			{
				ParticipantID: uuid.MustParse("b8c4d3e5-2f7f-4d8a-9b6a-5f7d1e2a3c4b"),
				Timestamp:     time.Now(),
			},
		},
		Read:      []entity.Event{}, // not yet read
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:        uuid.MustParse("66666666-7777-8888-9999-000000000000"),
		MessageID: uuid.MustParse("66666666-7777-8888-9999-000000000000"),
		Delivered: []entity.Event{
			{
				ParticipantID: uuid.MustParse("a7d3b2c4-3e6e-4b57-9c8d-3d9d1efae3c5"),
				Timestamp:     time.Now(),
			},
		},
		Read: []entity.Event{
			{
				ParticipantID: uuid.MustParse("a7d3b2c4-3e6e-4b57-9c8d-3d9d1efae3c5"),
				Timestamp:     time.Now(),
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
}
