package twitch_test

import (
	"KitsuneSemCalda/KitsuneBot/internal/twitch"
	"testing"

	"github.com/caiolandgraf/gest/v2/gest"
)

func TestNewClient(t *testing.T) {
	gest.Describe("NewClient").
		It("should create client with config", func(t *gest.T) {
			config := &twitch.Config{
				Username:   "testbot",
				OAuthToken: "oauth:test",
				Channel:    "testchannel",
			}

			client := twitch.NewClient(config)

			t.Expect(client).Not().ToBeNil()
		}).
		It("should create client with default config", func(t *gest.T) {
			config := twitch.LoadConfig()

			client := twitch.NewClient(config)

			t.Expect(client).Not().ToBeNil()
		}).
		Run(t)
}

func TestClientState(t *testing.T) {
	gest.Describe("Client - State Management").
		It("should start as not connected", func(t *gest.T) {
			config := twitch.LoadConfig()
			client := twitch.NewClient(config)

			t.Expect(client.IsConnected()).ToEqual(false)
		}).
		Run(t)
}

func TestClientRegisterHandlers(t *testing.T) {
	gest.Describe("Client - Handler Registration").
		It("should register single handler", func(t *gest.T) {
			config := twitch.LoadConfig()
			client := twitch.NewClient(config)
			registry := twitch.NewHandlerRegistry()
			registry.Register(twitch.NewPingHandler())

			client.RegisterHandlers(registry)

			t.Expect(registry.Count()).ToEqual(1)
		}).
		It("should register multiple handlers", func(t *gest.T) {
			config := twitch.LoadConfig()
			client := twitch.NewClient(config)
			registry := twitch.NewHandlerRegistry()
			registry.Register(twitch.NewPingHandler())
			registry.Register(twitch.NewStatsHandler())
			registry.Register(twitch.NewHelpHandler(registry))

			client.RegisterHandlers(registry)

			t.Expect(registry.Count()).ToEqual(3)
		}).
		Run(t)
}
