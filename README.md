# 🧠 Motion Dashboard

A minimalist personal dashboard to track your current and upcoming tasks, synced with Motion's API. Designed especially for neurodivergent minds (particularly ADHD), it provides visual timers, countdowns, and a calming way to stay oriented in time.

---

## 💡 Purpose

This was built to counter my own time-blindness. It's designed to sit quietly on a second screen or stay pinned in a browser tab, gently reminding me what I'm doing and what’s coming next—without yelling.

---

## 🔧 Features

- ⏰ Visual circle timer for the current task
- 🧾 Digital clock fallback when no task is active
- 📅 Countdown to the next scheduled task
- 🌙 Dark mode toggle with persistent preference via `localStorage`
- 🔁 Auto-refreshes data every 60 seconds
- 🔃 Smooth UI updates every second
- 🚫 Graceful fallback when no data is available

---

## 📦 Tech Stack

**Frontend:**

- Vanilla JavaScript
- HTML/CSS

**Backend (Go):**

- Graceful shutdown with context
- Periodic data polling from Motion API
- HTML rendering with `html/template`
- Structured logging via `log/slog`

---

## 🚀 Getting Started

### 🛠 Requirements

- Go 1.21 or higher
- A Motion API key
- A basic understanding of what day it is (optional, but recommended)

### 1. Clone the Repository
```
git clone https://github.com/andrewdavis506/motion-dashboard.git
cd motion-dashboard
```
### 2. Set Up Your Environment
Copy the example environment file:
```
cp .env.example .env
```
Then edit .env with your actual values:
```
MOTION_API_KEY=your-api-key-here
PORT=8080
REFRESH_INTERVAL=60s
LOG_API_KEYS=false
```
Or export variables manually:
```
export MOTION_API_KEY=your-api-key-here
```
### 3. Run the Server
```
go run .
```
Then open: http://localhost:8080

⚙️ Configuration

Configuration values are prioritized in the following order:
1. Command-line flags (e.g. --port, --api-key)
2. Environment variables (.env or shell)
3. Sensible defaults (port 8080, refresh interval 60s, etc.)
 
Available .env Settings

| Variable           | Description                              |
|--------------------|------------------------------------------|
| `MOTION_API_KEY`   | Your Motion API key (required)          |
| `PORT`             | Port to run the server on (default: 8080)|
| `REFRESH_INTERVAL` | How often to poll the Motion API (default: 60s)|
| `LOG_API_KEYS`     | Log API keys for debugging (default: false)|
| `TEMPLATES_DIR`    | Directory for HTML templates (default: ./templates)|
| `STATIC_DIR`       | Directory for static files (default: ./static)|
---

## ✨ Dashboard Behavior

- Pulls fresh task data from Motion every 60 seconds
- UI updates every second for visual continuity
- Auto-refreshes the view if a task ends or a new one begins

---

## Project Structure

```
task-dashboard/
├── cmd/
│   └── server/        # Application entry point
│       └── main.go
├── config/           # Configuration management
├── internal/
│   ├── api/          # Motion API client
│   ├── models/       # Data structures
│   ├── service/      # Business logic
│   └── web/          # Web server
├── templates/        # HTML templates
├── static/           # CSS, JS files
├── go.mod            # Go module file
└── Makefile          # Build commands
```
---

## 👀 Screenshots

> Coming soon. Probably. 

---

## 🧠 Why I Made This

Couldn't find a dashboard that suited my needs and my love of Motion.

---

## 📚 What I Learned
**Everything. Every single thing.**

This was my first time building a project end-to-end using Go on the backend and vanilla JavaScript for a living, breathing frontend. Along the way, I learned:

- How to structure a Go web server
- Managing context in long-running Go services
- log/slog

---

## 🧼 Potential TODOs

- Raspberry Pi for Office Signage

--- 

🤘 Contributing

This started as a personal project, but if it speaks to you and you’ve got something to add, cool.

1. Fork it
2. Branch it
3. Build it
4. PR it

Just don’t break what’s already working unless it’s broken. In that case, thanks.us. This project was built for learning—and that includes you.

---

Made with ✨ and a lot of "Wait, what was I doing again?"

— Andrew
