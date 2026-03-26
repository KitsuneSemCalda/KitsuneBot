package twitch

import (
	"KitsuneSemCalda/KitsuneBot/internal/db"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/adeithe/go-twitch/irc"
)

type Client struct {
	config         *Config
	conn           *irc.Conn
	handlers       *HandlerRegistry
	connected      bool
	mu             sync.RWMutex
	disconnectChan chan struct{}
}

func NewClient(config *Config) *Client {
	client := &Client{
		config:         config,
		conn:           &irc.Conn{},
		disconnectChan: make(chan struct{}, 1),
	}

	client.setupCallbacks()
	return client
}

func (c *Client) setupCallbacks() {
	c.conn.OnMessage(func(msg irc.ChatMessage) {
		c.handleMessage(msg)
	})

	c.conn.OnDisconnect(func() {
		log.Printf("Disconnected from Twitch IRC")
		c.mu.Lock()
		if c.connected {
			c.connected = false
			select {
			case c.disconnectChan <- struct{}{}:
			default:
			}
		}
		c.mu.Unlock()
	})

	c.conn.OnReconnect(func() {
		log.Printf("Reconnecting to Twitch IRC...")
	})
}

func (c *Client) handleMessage(msg irc.ChatMessage) {
	db.InsertMessage(db.DB, msg.Sender.Username, msg.Text)

	if c.handlers != nil {
		response, err := c.handlers.HandleMessage(msg)
		if err != nil {
			log.Printf("Handler error: %v", err)
			return
		}
		if response != "" {
			c.conn.Say(c.config.Channel, response)
		}
	}
}

func (c *Client) RegisterHandlers(registry *HandlerRegistry) {
	c.handlers = registry
}

func (c *Client) Connect() error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-sigChan:
			log.Printf("Shutting down...")
			c.mu.Lock()
			c.connected = false
			c.mu.Unlock()
			c.conn.Close()
			return nil
		default:
			err := c.connectWithBackoff()
			if err != nil {
				log.Printf("Connection failed: %v", err)
				continue
			}

			c.mu.Lock()
			c.connected = true
			c.mu.Unlock()

			log.Printf("Connected and joined channel: %s", c.config.Channel)

			select {
			case <-c.disconnectChan:
				log.Printf("Disconnected, will retry...")
			case <-sigChan:
				log.Printf("Shutting down...")
				c.conn.Close()
				return nil
			}
		}
	}
}

func (c *Client) connectWithBackoff() error {
	cfg := c.config.ReconnectConfig
	delay := cfg.InitialDelay
	attempt := 1

	for {
		log.Printf("Connecting (attempt %d)...", attempt)

		err := c.connectOnce()
		if err == nil {
			return nil
		}

		log.Printf("Connection failed: %v. Waiting %v before retry...", err, delay)
		time.Sleep(delay)

		delay = time.Duration(float64(delay) * cfg.Multiplier)
		if delay > cfg.MaxDelay {
			delay = cfg.MaxDelay
		}

		attempt++

		if !c.shouldContinue() {
			return fmt.Errorf("shutdown requested")
		}
	}
}

func (c *Client) connectOnce() error {
	err := c.conn.SetLogin(c.config.Username, c.config.OAuthToken)
	if err != nil {
		return err
	}

	err = c.conn.Connect()
	if err != nil {
		return err
	}

	err = c.conn.Join(c.config.Channel)
	if err != nil {
		return err
	}

	log.Printf("Connected to Twitch IRC as %s", c.config.Username)

	return nil
}

func (c *Client) shouldContinue() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return true
}

func (c *Client) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connected
}
