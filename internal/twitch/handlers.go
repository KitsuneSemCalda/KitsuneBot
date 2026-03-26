package twitch

import "github.com/adeithe/go-twitch/irc"

type Handler interface {
	Name() string
	Handle(msg irc.ChatMessage) (string, error)
	Usage() string
}

type HandlerRegistry struct {
	handlers map[string]Handler
}

func NewHandlerRegistry() *HandlerRegistry {
	return &HandlerRegistry{
		handlers: make(map[string]Handler),
	}
}

func (r *HandlerRegistry) Register(handler Handler) {
	r.handlers[handler.Name()] = handler
}

func (r *HandlerRegistry) Count() int {
	return len(r.handlers)
}

func (r *HandlerRegistry) HandleMessage(msg irc.ChatMessage) (string, error) {
	if len(msg.Text) == 0 || msg.Text[0] != '!' {
		return "", nil
	}
	parts := ParseCommand(msg.Text)
	if len(parts) == 0 {
		return "", nil
	}
	cmd := parts[0]
	handler, exists := r.handlers[cmd]
	if !exists {
		return "", nil
	}
	return handler.Handle(msg)
}

func ParseCommand(msg string) []string {
	var parts []string
	current := ""
	for _, ch := range msg {
		if ch == ' ' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}
