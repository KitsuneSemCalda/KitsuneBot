package twitch

import (
	"KitsuneSemCalda/KitsuneBot/internal/db"
	"fmt"

	irc "github.com/adeithe/go-twitch/irc"
)

type StatsHandler struct{}

func NewStatsHandler() *StatsHandler {
	return &StatsHandler{}
}

func (h *StatsHandler) Name() string {
	return "!stats"
}

func (h *StatsHandler) Usage() string {
	return "!stats - Mostra estatísticas do chat"
}

func (h *StatsHandler) Handle(msg irc.ChatMessage) (string, error) {
	count, err := db.GetMessageCount(db.DB)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Total de mensagens no banco: %d", count), nil
}
