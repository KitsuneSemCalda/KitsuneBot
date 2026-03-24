# KitsuneBot

A simple yet evolving Twitch bot written in Go, combining smart moderation with chat-driven behavior and emergent interactions.

---

## Features

- **Moderation:** Simple smart moderation using Bayesian classification.
- **Chat Interaction:** Collects and processes chat messages.
- **Extensible:** Easily add new features and behaviors.
- **Cross-platform:** Works on Windows, Linux, and macOS.
- **Tested:** Uses Go's native testing plus optional Jest-inspired `gest` for readable output.

---

## Project Structure

```
├── 📁 assets
├── 📁 cmd
│   └── 📁 kitsune_bot
│       └── 🐹 main.go
├── ⚙️ .gitignore
├── 📄 LICENSE
├── 📄 Makefile
├── 📝 README.md
├── 📄 go.mod
└── 📄 go.sum
```


---

## Getting Started

### Build

>[!WARNING]
> Check if CGO_ENABLE=1 before building, this is a necessary dependencie from functional database

```bash
make build
``` 

### Run

```bash
make run
```

### Tests

```bash
make test        # native Go tests
make test-pretty # optional, uses gest for nicer output
``` 

### Clean

```bash
make clean
``` 

