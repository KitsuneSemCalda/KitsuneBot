package twitch

import (
	"fmt"

	irc "github.com/adeithe/go-twitch/irc"
)

type PingHandler struct{}

func NewPingHandler() *PingHandler {
	return &PingHandler{}
}

func (h *PingHandler) Name() string {
	return "!ping"
}

func (h *PingHandler) Usage() string {
	return "!ping - Responde com pong"
}

func (h *PingHandler) Handle(msg irc.ChatMessage) (string, error) {
	return fmt.Sprintf("Pong! @%s", msg.Sender.Username), nil
}
