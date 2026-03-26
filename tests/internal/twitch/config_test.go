package twitch_test

import (
	"KitsuneSemCalda/KitsuneBot/internal/twitch"
	"os"
	"testing"
	"time"

	"github.com/caiolandgraf/gest/v2/gest"
)

func TestLoadConfig(t *testing.T) {
	gest.Describe("LoadConfig - Default Values").
		It("should load default values when env vars are not set", func(t *gest.T) {
			os.Unsetenv("TWITCH_USERNAME")
			os.Unsetenv("TWITCH_OAUTH")
			os.Unsetenv("TWITCH_CLIENT_ID")
			os.Unsetenv("TWITCH_CHANNEL")
			os.Unsetenv("DB_PATH")

			config := twitch.LoadConfig()

			t.Expect(config.Username).ToEqual("kitsunebot")
			t.Expect(config.OAuthToken).ToEqual("oauth:placeholder")
			t.Expect(config.ClientID).ToEqual("client_id_placeholder")
			t.Expect(config.Channel).ToEqual("kitsunebot")
			t.Expect(config.DbPath).ToEqual("./kitsunebot.db")
		}).
		It("should load default reconnect config when env vars are not set", func(t *gest.T) {
			os.Unsetenv("RECONNECT_INITIAL_DELAY")
			os.Unsetenv("RECONNECT_MAX_DELAY")
			os.Unsetenv("RECONNECT_MULTIPLIER")

			config := twitch.LoadConfig()

			t.Expect(config.ReconnectConfig.InitialDelay).ToEqual(1 * time.Second)
			t.Expect(config.ReconnectConfig.MaxDelay).ToEqual(30 * time.Second)
			t.Expect(config.ReconnectConfig.Multiplier).ToEqual(2.0)
		}).
		Run(t)
}

func TestLoadConfigFromEnv(t *testing.T) {
	gest.Describe("LoadConfig - Environment Variables").
		It("should load values from environment variables", func(t *gest.T) {
			os.Setenv("TWITCH_USERNAME", "testbot")
			os.Setenv("TWITCH_OAUTH", "oauth:testtoken")
			os.Setenv("TWITCH_CLIENT_ID", "testclientid")
			os.Setenv("TWITCH_CHANNEL", "testchannel")
			os.Setenv("DB_PATH", "./test.db")
			defer func() {
				os.Unsetenv("TWITCH_USERNAME")
				os.Unsetenv("TWITCH_OAUTH")
				os.Unsetenv("TWITCH_CLIENT_ID")
				os.Unsetenv("TWITCH_CHANNEL")
				os.Unsetenv("DB_PATH")
			}()

			config := twitch.LoadConfig()

			t.Expect(config.Username).ToEqual("testbot")
			t.Expect(config.OAuthToken).ToEqual("oauth:testtoken")
			t.Expect(config.ClientID).ToEqual("testclientid")
			t.Expect(config.Channel).ToEqual("testchannel")
			t.Expect(config.DbPath).ToEqual("./test.db")
		}).
		It("should load reconnect config from environment variables", func(t *gest.T) {
			os.Setenv("RECONNECT_INITIAL_DELAY", "2s")
			os.Setenv("RECONNECT_MAX_DELAY", "60s")
			os.Setenv("RECONNECT_MULTIPLIER", "3")
			defer func() {
				os.Unsetenv("RECONNECT_INITIAL_DELAY")
				os.Unsetenv("RECONNECT_MAX_DELAY")
				os.Unsetenv("RECONNECT_MULTIPLIER")
			}()

			config := twitch.LoadConfig()

			t.Expect(config.ReconnectConfig.InitialDelay).ToEqual(2 * time.Second)
			t.Expect(config.ReconnectConfig.MaxDelay).ToEqual(60 * time.Second)
			t.Expect(config.ReconnectConfig.Multiplier).ToEqual(3.0)
		}).
		Run(t)
}

func TestLoadConfigEdgeCases(t *testing.T) {
	gest.Describe("LoadConfig - Edge Cases").
		It("should handle empty string env vars as not set", func(t *gest.T) {
			os.Setenv("TWITCH_USERNAME", "")
			defer os.Unsetenv("TWITCH_USERNAME")

			config := twitch.LoadConfig()

			t.Expect(config.Username).ToEqual("kitsunebot")
		}).
		It("should handle invalid duration gracefully", func(t *gest.T) {
			os.Setenv("RECONNECT_INITIAL_DELAY", "invalid")
			defer os.Unsetenv("RECONNECT_INITIAL_DELAY")

			config := twitch.LoadConfig()

			t.Expect(config.ReconnectConfig.InitialDelay).ToEqual(1 * time.Second)
		}).
		It("should handle invalid max delay gracefully", func(t *gest.T) {
			os.Setenv("RECONNECT_MAX_DELAY", "invalid")
			defer os.Unsetenv("RECONNECT_MAX_DELAY")

			config := twitch.LoadConfig()

			t.Expect(config.ReconnectConfig.MaxDelay).ToEqual(30 * time.Second)
		}).
		It("should handle invalid multiplier gracefully", func(t *gest.T) {
			os.Setenv("RECONNECT_MULTIPLIER", "invalid")
			defer os.Unsetenv("RECONNECT_MULTIPLIER")

			config := twitch.LoadConfig()

			t.Expect(config.ReconnectConfig.Multiplier).ToEqual(2.0)
		}).
		It("should handle zero multiplier gracefully", func(t *gest.T) {
			os.Setenv("RECONNECT_MULTIPLIER", "0")
			defer os.Unsetenv("RECONNECT_MULTIPLIER")

			config := twitch.LoadConfig()

			t.Expect(config.ReconnectConfig.Multiplier).ToEqual(2.0)
		}).
		It("should handle negative multiplier gracefully", func(t *gest.T) {
			os.Setenv("RECONNECT_MULTIPLIER", "-1")
			defer os.Unsetenv("RECONNECT_MULTIPLIER")

			config := twitch.LoadConfig()

			t.Expect(config.ReconnectConfig.Multiplier).ToEqual(2.0)
		}).
		It("should handle very small initial delay", func(t *gest.T) {
			os.Setenv("RECONNECT_INITIAL_DELAY", "1ms")
			defer os.Unsetenv("RECONNECT_INITIAL_DELAY")

			config := twitch.LoadConfig()

			t.Expect(config.ReconnectConfig.InitialDelay).ToEqual(1 * time.Millisecond)
		}).
		It("should handle very large initial delay", func(t *gest.T) {
			os.Setenv("RECONNECT_INITIAL_DELAY", "1h")
			defer os.Unsetenv("RECONNECT_INITIAL_DELAY")

			config := twitch.LoadConfig()

			t.Expect(config.ReconnectConfig.InitialDelay).ToEqual(1 * time.Hour)
		}).
		It("should handle zero initial delay gracefully", func(t *gest.T) {
			os.Setenv("RECONNECT_INITIAL_DELAY", "0")
			defer os.Unsetenv("RECONNECT_INITIAL_DELAY")

			config := twitch.LoadConfig()

			t.Expect(config.ReconnectConfig.InitialDelay).ToEqual(1 * time.Second)
		}).
		It("should handle max delay smaller than initial delay", func(t *gest.T) {
			os.Setenv("RECONNECT_INITIAL_DELAY", "10s")
			os.Setenv("RECONNECT_MAX_DELAY", "1s")
			defer func() {
				os.Unsetenv("RECONNECT_INITIAL_DELAY")
				os.Unsetenv("RECONNECT_MAX_DELAY")
			}()

			config := twitch.LoadConfig()

			t.Expect(config.ReconnectConfig.InitialDelay).ToEqual(10 * time.Second)
			t.Expect(config.ReconnectConfig.MaxDelay).ToEqual(1 * time.Second)
		}).
		It("should handle special characters in channel name", func(t *gest.T) {
			os.Setenv("TWITCH_CHANNEL", "test_channel_123")
			defer os.Unsetenv("TWITCH_CHANNEL")

			config := twitch.LoadConfig()

			t.Expect(config.Channel).ToEqual("test_channel_123")
		}).
		It("should handle oauth token without prefix", func(t *gest.T) {
			os.Setenv("TWITCH_OAUTH", "abcdef123456")
			defer os.Unsetenv("TWITCH_OAUTH")

			config := twitch.LoadConfig()

			t.Expect(config.OAuthToken).ToEqual("abcdef123456")
		}).
		It("should preserve original case of username", func(t *gest.T) {
			os.Setenv("TWITCH_USERNAME", "TestBot")
			defer os.Unsetenv("TWITCH_USERNAME")

			config := twitch.LoadConfig()

			t.Expect(config.Username).ToEqual("TestBot")
		}).
		Run(t)
}
