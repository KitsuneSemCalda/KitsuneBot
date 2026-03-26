package twitch

import (
	"fmt"

	irc "github.com/adeithe/go-twitch/irc"
)

type HelpHandler struct {
	registry *HandlerRegistry
}

func NewHelpHandler(registry *HandlerRegistry) *HelpHandler {
	return &HelpHandler{registry: registry}
}

func (h *HelpHandler) Name() string {
	return "!help"
}

func (h *HelpHandler) Usage() string {
	return "!help - Lista todos os comandos"
}

func (h *HelpHandler) Handle(msg irc.ChatMessage) (string, error) {
	var commands []string
	for name, handler := range h.registry.handlers {
		commands = append(commands, fmt.Sprintf("%s: %s", name, handler.Usage()))
	}
	return "Comandos: " + joinStrings(commands, " | "), nil
}

func joinStrings(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}
