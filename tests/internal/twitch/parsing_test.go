package twitch_test

import (
	"KitsuneSemCalda/KitsuneBot/internal/twitch"
	"testing"

	"github.com/caiolandgraf/gest/v2/gest"
)

func TestParseCommand(t *testing.T) {
	gest.Describe("ParseCommand - Basic Cases").
		It("should parse simple command", func(t *gest.T) {
			result := twitch.ParseCommand("!ping")

			t.Expect(result).ToHaveLength(1)
			t.Expect(result[0]).ToEqual("!ping")
		}).
		It("should parse command with single argument", func(t *gest.T) {
			result := twitch.ParseCommand("!help ping")

			t.Expect(result).ToHaveLength(2)
			t.Expect(result[0]).ToEqual("!help")
			t.Expect(result[1]).ToEqual("ping")
		}).
		It("should parse command with multiple arguments", func(t *gest.T) {
			result := twitch.ParseCommand("!command arg1 arg2 arg3")

			t.Expect(result).ToHaveLength(4)
			t.Expect(result[0]).ToEqual("!command")
			t.Expect(result[1]).ToEqual("arg1")
			t.Expect(result[2]).ToEqual("arg2")
			t.Expect(result[3]).ToEqual("arg3")
		}).
		Run(t)
}

func TestParseCommandEdgeCases(t *testing.T) {
	gest.Describe("ParseCommand - Edge Cases").
		It("should handle empty string", func(t *gest.T) {
			result := twitch.ParseCommand("")

			t.Expect(result).ToHaveLength(0)
		}).
		It("should handle single space", func(t *gest.T) {
			result := twitch.ParseCommand(" ")

			t.Expect(result).ToHaveLength(0)
		}).
		It("should handle multiple spaces", func(t *gest.T) {
			result := twitch.ParseCommand("   ")

			t.Expect(result).ToHaveLength(0)
		}).
		It("should handle leading spaces", func(t *gest.T) {
			result := twitch.ParseCommand("  !ping")

			t.Expect(result).ToHaveLength(1)
			t.Expect(result[0]).ToEqual("!ping")
		}).
		It("should handle trailing spaces", func(t *gest.T) {
			result := twitch.ParseCommand("!ping  ")

			t.Expect(result).ToHaveLength(1)
			t.Expect(result[0]).ToEqual("!ping")
		}).
		It("should handle multiple consecutive spaces", func(t *gest.T) {
			result := twitch.ParseCommand("!ping   arg1    arg2")

			t.Expect(result).ToHaveLength(3)
			t.Expect(result[0]).ToEqual("!ping")
			t.Expect(result[1]).ToEqual("arg1")
			t.Expect(result[2]).ToEqual("arg2")
		}).
		It("should handle command with no arguments", func(t *gest.T) {
			result := twitch.ParseCommand("!ping")

			t.Expect(result).ToHaveLength(1)
			t.Expect(result[0]).ToEqual("!ping")
		}).
		It("should handle single character argument", func(t *gest.T) {
			result := twitch.ParseCommand("!ping a")

			t.Expect(result).ToHaveLength(2)
			t.Expect(result[1]).ToEqual("a")
		}).
		Run(t)
}

func TestParseCommandSpecialCases(t *testing.T) {
	gest.Describe("ParseCommand - Special Cases").
		It("should handle command with empty arguments", func(t *gest.T) {
			result := twitch.ParseCommand("!ping ")

			t.Expect(result).ToHaveLength(1)
			t.Expect(result[0]).ToEqual("!ping")
		}).
		It("should handle command with special characters", func(t *gest.T) {
			result := twitch.ParseCommand("!test arg1! arg2@ arg3#")

			t.Expect(result).ToHaveLength(4)
			t.Expect(result[1]).ToEqual("arg1!")
			t.Expect(result[2]).ToEqual("arg2@")
			t.Expect(result[3]).ToEqual("arg3#")
		}).
		It("should handle command with numbers", func(t *gest.T) {
			result := twitch.ParseCommand("!test 123 456")

			t.Expect(result).ToHaveLength(3)
			t.Expect(result[1]).ToEqual("123")
			t.Expect(result[2]).ToEqual("456")
		}).
		It("should handle command with underscores", func(t *gest.T) {
			result := twitch.ParseCommand("!test user_name test_value")

			t.Expect(result).ToHaveLength(3)
			t.Expect(result[1]).ToEqual("user_name")
			t.Expect(result[2]).ToEqual("test_value")
		}).
		It("should handle command with hyphens", func(t *gest.T) {
			result := twitch.ParseCommand("!test user-name test-value")

			t.Expect(result).ToHaveLength(3)
			t.Expect(result[1]).ToEqual("user-name")
			t.Expect(result[2]).ToEqual("test-value")
		}).
		It("should handle command with unicode", func(t *gest.T) {
			result := twitch.ParseCommand("!test こんにちは 世界")

			t.Expect(result).ToHaveLength(3)
			t.Expect(result[1]).ToEqual("こんにちは")
			t.Expect(result[2]).ToEqual("世界")
		}).
		It("should handle command with emojis", func(t *gest.T) {
			result := twitch.ParseCommand("!test 🎮 🎯")

			t.Expect(result).ToHaveLength(3)
			t.Expect(result[1]).ToEqual("🎮")
			t.Expect(result[2]).ToEqual("🎯")
		}).
		It("should handle long arguments", func(t *gest.T) {
			longArg := "Lorem"
			result := twitch.ParseCommand("!test " + longArg)

			t.Expect(result).ToHaveLength(2)
			t.Expect(result[1]).ToEqual(longArg)
		}).
		Run(t)
}

func TestParseCommandRealWorld(t *testing.T) {
	gest.Describe("ParseCommand - Real World Examples").
		It("should parse !help command", func(t *gest.T) {
			result := twitch.ParseCommand("!help")

			t.Expect(result).ToHaveLength(1)
			t.Expect(result[0]).ToEqual("!help")
		}).
		It("should parse !ping command", func(t *gest.T) {
			result := twitch.ParseCommand("!ping")

			t.Expect(result).ToHaveLength(1)
			t.Expect(result[0]).ToEqual("!ping")
		}).
		It("should parse !stats command", func(t *gest.T) {
			result := twitch.ParseCommand("!stats")

			t.Expect(result).ToHaveLength(1)
			t.Expect(result[0]).ToEqual("!stats")
		}).
		It("should parse command with mention", func(t *gest.T) {
			result := twitch.ParseCommand("!hello @username")

			t.Expect(result).ToHaveLength(2)
			t.Expect(result[0]).ToEqual("!hello")
			t.Expect(result[1]).ToEqual("@username")
		}).
		It("should parse command with multiple mentions", func(t *gest.T) {
			result := twitch.ParseCommand("!notify @user1 @user2 @user3")

			t.Expect(result).ToHaveLength(4)
			t.Expect(result[0]).ToEqual("!notify")
			t.Expect(result[1]).ToEqual("@user1")
			t.Expect(result[2]).ToEqual("@user2")
			t.Expect(result[3]).ToEqual("@user3")
		}).
		It("should parse command with channel name", func(t *gest.T) {
			result := twitch.ParseCommand("!join #testchannel")

			t.Expect(result).ToHaveLength(2)
			t.Expect(result[0]).ToEqual("!join")
			t.Expect(result[1]).ToEqual("#testchannel")
		}).
		It("should parse command with number argument", func(t *gest.T) {
			result := twitch.ParseCommand("!timeout user123 600")

			t.Expect(result).ToHaveLength(3)
			t.Expect(result[0]).ToEqual("!timeout")
			t.Expect(result[1]).ToEqual("user123")
			t.Expect(result[2]).ToEqual("600")
		}).
		It("should parse command with reason", func(t *gest.T) {
			result := twitch.ParseCommand("!ban user123 spamming links")

			t.Expect(result).ToHaveLength(4)
			t.Expect(result[0]).ToEqual("!ban")
			t.Expect(result[1]).ToEqual("user123")
			t.Expect(result[2]).ToEqual("spamming")
			t.Expect(result[3]).ToEqual("links")
		}).
		Run(t)
}
