#!/bin/bash

# SpecForge Development Environment Startup Script (Bash/Git Bash/WSL)
# Mirrors logic from start_dev.ps1 for cross-platform availability.

# ----------------------------------------------------
# 0. Configuration & Colors
# ----------------------------------------------------
ONLY_MIGRATIONS=false

for arg in "$@"; do
    case $arg in
        --only-migrations)
        ONLY_MIGRATIONS=true
        shift
        ;;
    esac
done

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
CYAN='\033[0;36m'
GRAY='\033[0;90m'
NC='\033[0m' # No Color

echo -e "\n${CYAN}====================================================${NC}"
echo -e "${CYAN}          SpecForge Dev Environment Start          ${NC}"
echo -e "${CYAN}====================================================${NC}\n"

# ----------------------------------------------------
# 1. Prerequisite Checks
# ----------------------------------------------------
echo -e "${YELLOW}0. Checking Prerequisites...${NC}"

# Docker
if ! command -v docker &> /dev/null; then
    echo -e "${RED} ❌ Docker is not installed or not in PATH!${NC}"
    exit 1
fi

# Go
if ! command -v go &> /dev/null; then
    echo -e "${RED} ❌ Go is not installed! (Required for Backend)${NC}"
    exit 1
fi

# Node/NPM
if ! command -v npm &> /dev/null; then
    echo -e "${RED} ❌ NPM is not installed! (Required for Frontend)${NC}"
    exit 1
fi

echo -e "${GREEN} ✅ All prerequisites met${NC}"

# ----------------------------------------------------
# 2. Start Infrastructure (Docker)
# ----------------------------------------------------
echo -e "\n${YELLOW}1. Starting Database & Migrations (Docker)...${NC}"

docker compose up -d postgres migrate
if [ $? -eq 0 ]; then
    echo -e "${GREEN} ✅ Docker services started${NC}"
else
    echo -e "${RED} ❌ Failed to start Docker services. Check docker-compose.yml${NC}"
    exit 1
fi

# Wait for Postgres Port 5432
echo -ne "${GRAY}   Waiting for Postgres (5432)...${NC}"
retries=30
while [ $retries -gt 0 ]; do
    if (echo > /dev/tcp/localhost/5432) >/dev/null 2>&1; then
        echo -e "${GREEN} Ready!${NC}"
        break
    fi
    echo -n "."
    sleep 1
    ((retries--))
done

if [ $retries -le 0 ]; then
    echo -e "${RED} ❌ Timeout waiting for Postgres.${NC}"
    exit 1
fi

if [ "$ONLY_MIGRATIONS" = true ]; then
    echo -e "\n${GREEN}🎉 Migrations completed successfully!${NC}"
    exit 0
fi

# ----------------------------------------------------
# 3. Start Application (Local)
# ----------------------------------------------------
echo -e "\n${YELLOW}2. Starting Backend (Go)...${NC}"

# Check if Backend is already running
if (echo > /dev/tcp/localhost/8080) >/dev/null 2>&1; then
    echo -e "${CYAN}   ℹ️  Backend appears to be running already on port 8080${NC}"
else
    echo -e "${GREEN}   🚀 Launching Backend in new window...${NC}"
    
    # Environment variables for local run (connecting to localhost DB)
    export DATABASE_URL='postgres://specforge:specforge@localhost:5432/specforge?sslmode=disable'
    export JWT_SECRET='supersecret'
    export PORT='8080'
    
    # Launch logic based on OS/Terminal
    if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" ]]; then
        # Windows (Git Bash) - Using 'start' to spawn a new CMD/Bash window
        start "SpecForge Backend" bash -c "echo 'Starting SpecForge Backend...'; cd backend && go run ./cmd/server; echo 'Process exited. Press Enter to close...'; read"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        osascript -e 'tell app "Terminal" to do script "cd '$(pwd)'/backend && export DATABASE_URL='\''postgres://specforge:specforge@localhost:5432/specforge?sslmode=disable'\'' && export JWT_SECRET='\''supersecret'\'' && export PORT='\''8080'\'' && go run ./cmd/server"' 
    else
        # Linux (attempt common terminal emulators)
        if command -v gnome-terminal &> /dev/null; then
            gnome-terminal -- bash -c "cd backend && go run ./cmd/server; exec bash"
        elif command -v xterm &> /dev/null; then
            xterm -e "cd backend && go run ./cmd/server" &
        else
            echo -e "${YELLOW}   ⚠️  No terminal emulator found. Running in background (logs to backend.log)${NC}"
            (cd backend && go run ./cmd/server) > backend.log 2>&1 &
        fi
    fi
fi

echo -e "\n${YELLOW}3. Starting Frontend (Vite)...${NC}"

# Check if Frontend is already running
if (echo > /dev/tcp/localhost/5173) >/dev/null 2>&1; then
    echo -e "${CYAN}   ℹ️  Frontend appears to be running on port 5173${NC}"
else
    echo -e "${GREEN}   🚀 Launching Frontend in new window...${NC}"
    
    if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" ]]; then
        # Windows (Git Bash)
        start "SpecForge Frontend" bash -c "echo 'Starting SpecForge Frontend...'; cd frontend && npm run dev; echo 'Process exited. Press Enter to close...'; read"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        osascript -e 'tell app "Terminal" to do script "cd '$(pwd)'/frontend && npm run dev"'
    else
        # Linux
        if command -v gnome-terminal &> /dev/null; then
            gnome-terminal -- bash -c "cd frontend && npm run dev; exec bash"
        elif command -v xterm &> /dev/null; then
            xterm -e "cd frontend && npm run dev" &
        else
            echo -e "${YELLOW}   ⚠️  No terminal emulator found. Running in background (logs to frontend.log)${NC}"
            (cd frontend && npm run dev) > frontend.log 2>&1 &
        fi
    fi
fi

# ----------------------------------------------------
# 4. Verification
# ----------------------------------------------------
echo -e "\n${YELLOW}4. Verifying Endpoints (Waiting for startup)...${NC}"

# Poll Backend Health
backendReady=false
echo -ne "${GRAY}   Waiting for Backend...${NC}"
for ((i=0; i<20; i++)); do
    if curl -s http://localhost:8080/health | grep -q "OK" &> /dev/null; then
        backendReady=true
        break
    fi
    echo -n "."
    sleep 1
done
echo ""

if [ "$backendReady" = true ]; then
    echo -e "${GREEN}    ✅ Backend is UP${NC}"
else
    echo -e "${RED}    ❌ Backend failed to start within timeout${NC}"
fi

# Poll Frontend
frontendReady=false
echo -ne "${GRAY}   Waiting for Frontend...${NC}"
for ((i=0; i<20; i++)); do
    if (echo > /dev/tcp/localhost/5173) >/dev/null 2>&1; then
        frontendReady=true
        break
    fi
    echo -n "."
    sleep 1
done
echo ""

if [ "$frontendReady" = true ]; then
    echo -e "${GREEN}    ✅ Frontend is UP${NC}"
else
    echo -e "${RED}    ❌ Frontend failed to start within timeout${NC}"
fi

# ----------------------------------------------------
# 5. Final Status
# ----------------------------------------------------
echo -e "\n${CYAN}----------------------------------------${NC}"
if [ "$backendReady" = true ] && [ "$frontendReady" = true ]; then
    echo -e "${GREEN}🎉 Development Environment is READY!${NC}"
    echo -e "   Frontend: ${CYAN}http://localhost:5173${NC}"
    echo -e "   Backend:  ${CYAN}http://localhost:8080${NC}"
else
    echo -e "${YELLOW} ⚠️  Some services failed to start correctly.${NC}"
    echo -e " Check windows/logs for details."
fi
echo -e "${CYAN}----------------------------------------${NC}\n"
