package twitch_test

import (
	"KitsuneSemCalda/KitsuneBot/internal/twitch"
	"testing"

	"github.com/adeithe/go-twitch/irc"
	"github.com/caiolandgraf/gest/v2/gest"
)

func TestHandlerRegistry(t *testing.T) {
	gest.Describe("HandlerRegistry - Creation").
		It("should create empty registry", func(t *gest.T) {
			registry := twitch.NewHandlerRegistry()

			t.Expect(registry).Not().ToBeNil()
			t.Expect(registry.Count()).ToEqual(0)
		}).
		Run(t)
}

func TestHandlerRegistryRegister(t *testing.T) {
	gest.Describe("HandlerRegistry - Register").
		It("should register handler", func(t *gest.T) {
			registry := twitch.NewHandlerRegistry()
			handler := twitch.NewPingHandler()

			registry.Register(handler)

			t.Expect(registry.Count()).ToEqual(1)
		}).
		It("should register multiple handlers", func(t *gest.T) {
			registry := twitch.NewHandlerRegistry()

			registry.Register(twitch.NewPingHandler())
			registry.Register(twitch.NewStatsHandler())
			registry.Register(twitch.NewHelpHandler(registry))

			t.Expect(registry.Count()).ToEqual(3)
		}).
		It("should overwrite existing handler with same name", func(t *gest.T) {
			registry := twitch.NewHandlerRegistry()

			registry.Register(twitch.NewPingHandler())
			registry.Register(twitch.NewPingHandler())

			t.Expect(registry.Count()).ToEqual(1)
		}).
		It("should allow handlers with different names", func(t *gest.T) {
			registry := twitch.NewHandlerRegistry()

			registry.Register(twitch.NewPingHandler())
			registry.Register(twitch.NewStatsHandler())

			t.Expect(registry.Count()).ToEqual(2)
		}).
		Run(t)
}

func TestHandlerRegistryHandleMessage(t *testing.T) {
	gest.Describe("HandlerRegistry - HandleMessage").
		It("should return empty for non-command messages", func(t *gest.T) {
			registry := twitch.NewHandlerRegistry()
			registry.Register(twitch.NewPingHandler())

			msg := irc.ChatMessage{
				Text: "Hello world",
				Sender: irc.ChatSender{
					Username: "testuser",
				},
			}

			response, err := registry.HandleMessage(msg)

			t.Expect(err).ToBeNil()
			t.Expect(response).ToEqual("")
		}).
		It("should return empty for unknown command", func(t *gest.T) {
			registry := twitch.NewHandlerRegistry()
			registry.Register(twitch.NewPingHandler())

			msg := irc.ChatMessage{
				Text: "!unknowncommand",
				Sender: irc.ChatSender{
					Username: "testuser",
				},
			}

			response, err := registry.HandleMessage(msg)

			t.Expect(err).ToBeNil()
			t.Expect(response).ToEqual("")
		}).
		It("should return empty for empty message", func(t *gest.T) {
			registry := twitch.NewHandlerRegistry()
			registry.Register(twitch.NewPingHandler())

			msg := irc.ChatMessage{
				Text: "",
				Sender: irc.ChatSender{
					Username: "testuser",
				},
			}

			response, err := registry.HandleMessage(msg)

			t.Expect(err).ToBeNil()
			t.Expect(response).ToEqual("")
		}).
		It("should handle command with arguments", func(t *gest.T) {
			registry := twitch.NewHandlerRegistry()
			registry.Register(twitch.NewPingHandler())

			msg := irc.ChatMessage{
				Text: "!ping arg1 arg2",
				Sender: irc.ChatSender{
					Username: "testuser",
				},
			}

			response, err := registry.HandleMessage(msg)

			t.Expect(err).ToBeNil()
			t.Expect(response).ToContain("Pong!")
		}).
		It("should handle command with trailing spaces", func(t *gest.T) {
			registry := twitch.NewHandlerRegistry()
			registry.Register(twitch.NewPingHandler())

			msg := irc.ChatMessage{
				Text: "!ping  ",
				Sender: irc.ChatSender{
					Username: "testuser",
				},
			}

			response, err := registry.HandleMessage(msg)

			t.Expect(err).ToBeNil()
			t.Expect(response).ToContain("Pong!")
		}).
		It("should return empty for message with only spaces", func(t *gest.T) {
			registry := twitch.NewHandlerRegistry()
			registry.Register(twitch.NewPingHandler())

			msg := irc.ChatMessage{
				Text: "   ",
				Sender: irc.ChatSender{
					Username: "testuser",
				},
			}

			response, err := registry.HandleMessage(msg)

			t.Expect(err).ToBeNil()
			t.Expect(response).ToEqual("")
		}).
		It("should not treat numbers as commands", func(t *gest.T) {
			registry := twitch.NewHandlerRegistry()
			registry.Register(twitch.NewPingHandler())

			msg := irc.ChatMessage{
				Text: "12345",
				Sender: irc.ChatSender{
					Username: "testuser",
				},
			}

			response, err := registry.HandleMessage(msg)

			t.Expect(err).ToBeNil()
			t.Expect(response).ToEqual("")
		}).
		It("should not treat special characters as commands", func(t *gest.T) {
			registry := twitch.NewHandlerRegistry()
			registry.Register(twitch.NewPingHandler())

			msg := irc.ChatMessage{
				Text: "!@#$%",
				Sender: irc.ChatSender{
					Username: "testuser",
				},
			}

			response, err := registry.HandleMessage(msg)

			t.Expect(err).ToBeNil()
			t.Expect(response).ToEqual("")
		}).
		Run(t)
}

func TestPingHandler(t *testing.T) {
	gest.Describe("PingHandler - Basic").
		It("should return correct name", func(t *gest.T) {
			handler := twitch.NewPingHandler()

			t.Expect(handler.Name()).ToEqual("!ping")
		}).
		It("should return correct usage", func(t *gest.T) {
			handler := twitch.NewPingHandler()

			t.Expect(handler.Usage()).ToEqual("!ping - Responde com pong")
		}).
		Run(t)
}

func TestPingHandlerResponses(t *testing.T) {
	gest.Describe("PingHandler - Responses").
		It("should respond with pong to user", func(t *gest.T) {
			handler := twitch.NewPingHandler()

			msg := irc.ChatMessage{
				Text: "!ping",
				Sender: irc.ChatSender{
					Username: "testuser",
				},
			}

			response, err := handler.Handle(msg)

			t.Expect(err).ToBeNil()
			t.Expect(response).ToContain("Pong!")
			t.Expect(response).ToContain("@testuser")
		}).
		It("should handle different usernames", func(t *gest.T) {
			handler := twitch.NewPingHandler()

			testCases := []string{"alice", "bob123", "user_name", "UPPERCASE", "测试用户"}

			for _, username := range testCases {
				msg := irc.ChatMessage{
					Text: "!ping",
					Sender: irc.ChatSender{
						Username: username,
					},
				}

				response, err := handler.Handle(msg)

				t.Expect(err).ToBeNil()
				t.Expect(response).ToContain("@" + username)
			}
		}).
		It("should handle message with extra arguments", func(t *gest.T) {
			handler := twitch.NewPingHandler()

			msg := irc.ChatMessage{
				Text: "!ping arg1 arg2 arg3",
				Sender: irc.ChatSender{
					Username: "testuser",
				},
			}

			response, err := handler.Handle(msg)

			t.Expect(err).ToBeNil()
			t.Expect(response).ToContain("Pong!")
		}).
		Run(t)
}

func TestStatsHandler(t *testing.T) {
	gest.Describe("StatsHandler - Basic").
		It("should return correct name", func(t *gest.T) {
			handler := twitch.NewStatsHandler()

			t.Expect(handler.Name()).ToEqual("!stats")
		}).
		It("should return correct usage", func(t *gest.T) {
			handler := twitch.NewStatsHandler()

			t.Expect(handler.Usage()).ToEqual("!stats - Mostra estatísticas do chat")
		}).
		Run(t)
}

func TestHelpHandler(t *testing.T) {
	gest.Describe("HelpHandler - Basic").
		It("should return correct name", func(t *gest.T) {
			registry := twitch.NewHandlerRegistry()
			handler := twitch.NewHelpHandler(registry)

			t.Expect(handler.Name()).ToEqual("!help")
		}).
		It("should return correct usage", func(t *gest.T) {
			registry := twitch.NewHandlerRegistry()
			handler := twitch.NewHelpHandler(registry)

			t.Expect(handler.Usage()).ToEqual("!help - Lista todos os comandos")
		}).
		It("should list all registered commands", func(t *gest.T) {
			registry := twitch.NewHandlerRegistry()
			registry.Register(twitch.NewPingHandler())
			registry.Register(twitch.NewStatsHandler())

			handler := twitch.NewHelpHandler(registry)

			msg := irc.ChatMessage{
				Text: "!help",
				Sender: irc.ChatSender{
					Username: "testuser",
				},
			}

			response, err := handler.Handle(msg)

			t.Expect(err).ToBeNil()
			t.Expect(response).ToContain("Comandos:")
			t.Expect(response).ToContain("!ping")
			t.Expect(response).ToContain("!stats")
		}).
		It("should handle empty registry gracefully", func(t *gest.T) {
			registry := twitch.NewHandlerRegistry()

			handler := twitch.NewHelpHandler(registry)

			msg := irc.ChatMessage{
				Text: "!help",
				Sender: irc.ChatSender{
					Username: "testuser",
				},
			}

			response, err := handler.Handle(msg)

			t.Expect(err).ToBeNil()
			t.Expect(response).ToContain("Comandos:")
		}).
		It("should list multiple commands in order", func(t *gest.T) {
			registry := twitch.NewHandlerRegistry()
			registry.Register(twitch.NewPingHandler())
			registry.Register(twitch.NewStatsHandler())
			registry.Register(twitch.NewHelpHandler(registry))

			handler := twitch.NewHelpHandler(registry)

			msg := irc.ChatMessage{
				Text: "!help",
				Sender: irc.ChatSender{
					Username: "testuser",
				},
			}

			response, err := handler.Handle(msg)

			t.Expect(err).ToBeNil()
			t.Expect(response).ToContain("!ping")
			t.Expect(response).ToContain("!stats")
			t.Expect(response).ToContain("!help")
		}).
		Run(t)
}

func TestHandlerMessageFields(t *testing.T) {
	gest.Describe("Handlers - Message Fields").
		It("should access message ID field", func(t *gest.T) {
			handler := twitch.NewPingHandler()

			msg := irc.ChatMessage{
				ID:   "msg-id-123",
				Text: "!ping",
				Sender: irc.ChatSender{
					Username: "testuser",
				},
			}

			response, err := handler.Handle(msg)

			t.Expect(err).ToBeNil()
			t.Expect(response).ToContain("Pong!")
		}).
		It("should access message channel field", func(t *gest.T) {
			handler := twitch.NewPingHandler()

			msg := irc.ChatMessage{
				Channel: "#testchannel",
				Text:    "!ping",
				Sender: irc.ChatSender{
					Username: "testuser",
				},
			}

			response, err := handler.Handle(msg)

			t.Expect(err).ToBeNil()
			t.Expect(response).ToContain("Pong!")
		}).
		It("should access sender display name", func(t *gest.T) {
			handler := twitch.NewPingHandler()

			msg := irc.ChatMessage{
				Text: "!ping",
				Sender: irc.ChatSender{
					Username:    "testuser",
					DisplayName: "TestUser",
				},
			}

			response, err := handler.Handle(msg)

			t.Expect(err).ToBeNil()
			t.Expect(response).ToContain("@testuser")
		}).
		It("should access sender ID", func(t *gest.T) {
			handler := twitch.NewPingHandler()

			msg := irc.ChatMessage{
				Text: "!ping",
				Sender: irc.ChatSender{
					Username: "testuser",
					ID:       12345,
				},
			}

			response, err := handler.Handle(msg)

			t.Expect(err).ToBeNil()
			t.Expect(response).ToContain("Pong!")
		}).
		It("should handle sender with badges", func(t *gest.T) {
			handler := twitch.NewPingHandler()

			msg := irc.ChatMessage{
				Text: "!ping",
				Sender: irc.ChatSender{
					Username:     "testuser",
					IsModerator:  true,
					IsSubscriber: true,
					IsVIP:        true,
					Badges: map[string]string{
						"moderator":  "1",
						"subscriber": "12",
					},
				},
			}

			response, err := handler.Handle(msg)

			t.Expect(err).ToBeNil()
			t.Expect(response).ToContain("Pong!")
		}).
		It("should handle sender with color", func(t *gest.T) {
			handler := twitch.NewPingHandler()

			msg := irc.ChatMessage{
				Text: "!ping",
				Sender: irc.ChatSender{
					Username: "testuser",
					Color:    "#FF0000",
				},
			}

			response, err := handler.Handle(msg)

			t.Expect(err).ToBeNil()
			t.Expect(response).ToContain("Pong!")
		}).
		Run(t)
}

func TestHandlerUnicode(t *testing.T) {
	gest.Describe("Handlers - Unicode Support").
		It("should handle unicode usernames", func(t *gest.T) {
			handler := twitch.NewPingHandler()

			testUsernames := []string{
				"テストユーザー",
				"مستخدم",
				"用户",
				"üser",
				"ñ usuário",
			}

			for _, username := range testUsernames {
				msg := irc.ChatMessage{
					Text: "!ping",
					Sender: irc.ChatSender{
						Username: username,
					},
				}

				response, err := handler.Handle(msg)

				t.Expect(err).ToBeNil()
				t.Expect(response).ToContain("@" + username)
			}
		}).
		It("should handle unicode message content", func(t *gest.T) {
			handler := twitch.NewPingHandler()

			msg := irc.ChatMessage{
				Text: "🎉 !ping 🎉",
				Sender: irc.ChatSender{
					Username: "testuser",
				},
			}

			response, err := handler.Handle(msg)

			t.Expect(err).ToBeNil()
			t.Expect(response).ToContain("Pong!")
		}).
		Run(t)
}
