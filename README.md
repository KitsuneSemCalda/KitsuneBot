# KitsuneBot

A Twitch bot written in Go, combining smart moderation with chat-driven behavior and emergent interactions.

---

## Features

- **Twitch IRC Integration** - Connect to Twitch chat via go-twitch library
- **Auto-Reconnect** - Exponential backoff reconnection (1s → 2s → 4s → ... → 30s)
- **Graceful Shutdown** - Handles SIGINT/SIGTERM signals
- **Moderation** - Bayesian classification for smart moderation (planned)
- **Markov Chain** - Entertainment based on nMarkov style (planned)
- **Extensible** - Handler-based command system
- **SQLite Database** - Persistent message storage with TTL cleanup
- **Tested** - Comprehensive test suite using `gest` framework
- **Cross-platform** - Works on Windows, Linux, and macOS

---

## Dependencies

- [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) - SQLite database
- [caiolandgraf/gest](https://github.com/caiolandgraf/gest) - Testing framework
- [adeithe/go-twitch](https://github.com/adeithe/go-twitch) - Twitch IRC client
- [joho/godotenv](https://github.com/joho/godotenv) - Environment variables

---

## Project Structure

```
KitsuneBot/
├── cmd/
│   └── kitsune_bot/
│       └── main.go              # Entry point
├── internal/
│   ├── db/                      # Database layer
│   │   ├── db.go               # SQLite operations
│   │   └── message.go          # Message struct
│   └── twitch/                  # Twitch integration
│       ├── config.go           # Configuration loading
│       ├── client.go           # IRC client with reconnection
│       ├── handlers.go         # Handler registry
│       ├── handler_ping.go     # !ping command
│       ├── handler_help.go     # !help command
│       └── handler_stats.go    # !stats command
├── tests/
│   └── internal/
│       ├── database/           # Database tests
│       └── twitch/             # Twitch tests
├── .env.example               # Environment template
├── go.mod
└── README.md
```

---

## Configuration

Create a `.env` file based on `.env.example`:

```env
# Required
TWITCH_USERNAME=kitsunebot
TWITCH_OAUTH=oauth:your_oauth_token_here
TWITCH_CLIENTID=your_client_id_here
TWITCH_CHANNEL=your_channel_name
DB_PATH=./bin/kitsunebot.db

# Optional - Reconnection settings
RECONNECT_INITIAL_DELAY=1s
RECONNECT_MAX_DELAY=30s
RECONNECT_MULTIPLIER=2
```

### Configuration Options

| Variable | Default | Description |
|----------|---------|-------------|
| `TWITCH_USERNAME` | `kitsunebot` | Bot username |
| `TWITCH_OAUTH` | `oauth:placeholder` | OAuth token (get from Twitch) |
| `TWITCH_CLIENT_ID` | `client_id_placeholder` | Twitch Client ID |
| `TWITCH_CHANNEL` | `kitsunebot` | Channel to join |
| `DB_PATH` | `./kitsunebot.db` | SQLite database path |
| `RECONNECT_INITIAL_DELAY` | `1s` | Initial delay between reconnection attempts |
| `RECONNECT_MAX_DELAY` | `30s` | Maximum delay between reconnection attempts |
| `RECONNECT_MULTIPLIER` | `2` | Multiplier for exponential backoff |

---

## Getting Started

### Prerequisites

- Go 1.26+
- CGO enabled (required for SQLite)

### Build

```bash
go build -o kitsunebot ./cmd/kitsune_bot/
```

Or use Make:

```bash
make build
```

### Run

```bash
./kitsunebot
```

Or use Make:

```bash
make run
```

### Tests

Run all tests:

```bash
go test ./tests/...
```

Or use Make:

```bash
make test        # native Go tests
make test-pretty # uses gest for nicer output
```

---

## Commands

The bot includes built-in commands:

| Command | Description |
|---------|-------------|
| `!ping` | Responds with pong |
| `!help` | Lists all available commands |
| `!stats` | Shows message count in database |

---

## Database

The bot uses SQLite to store messages with the following schema:

```sql
CREATE TABLE messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user TEXT NOT NULL,
    content TEXT NOT NULL,
    preprocessed_content TEXT,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Features

- **TTL Cleanup** - Messages older than 30 days are automatically deleted
- **Capacity Limit** - (Planned) Maximum message storage capacity
- **Message Compression** - (Planned) Logical compression of similar messages

---

## Architecture

### Client Connection Flow

```
Start
  │
  ▼
Connect ──► Success? ──► No ──► Wait (backoff) ──► Retry
  │              │
  │             Yes
  ▼              ▼
Connected    Success!
  │
  ▼
Wait for disconnect ──► Disconnected ──► Wait (backoff) ──► Connect
  ▲                                                    │
  │____________________________________________________|
```

### Exponential Backoff

```
Attempt 1: wait 1s
Attempt 2: wait 2s
Attempt 3: wait 4s
Attempt 4: wait 8s
Attempt 5: wait 16s
Attempt 6+: wait 30s (max)
```

---

## License

See LICENSE file.
