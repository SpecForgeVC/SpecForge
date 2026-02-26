param (
    [switch]$OnlyMigrations
)

# ----------------------------------------------------
# 0. Configuration
# ----------------------------------------------------
Write-Host "1. Starting Database & Migrations (Docker)..." -ForegroundColor Yellow
if (Get-Command "docker" -ErrorAction SilentlyContinue) {
    docker compose up -d postgres migrate
    if ($LASTEXITCODE -eq 0) {
        Write-Host " Docker services started" -ForegroundColor Green
    } else {
        Write-Host " Failed to start Docker services" -ForegroundColor Red
        exit 1
    }
} else {
    Write-Host " Docker is not installed!" -ForegroundColor Red
    exit 1
}

# Wait for DB port
Write-Host "   Waiting for Postgres (5432)..." -ForegroundColor DarkGray
$retries = 30
while ($retries -gt 0) {
    try {
        $tcp = New-Object System.Net.Sockets.TcpClient
        $tcp.Connect("localhost", 5432)
        $tcp.Close()
        break
    } catch {
        Start-Sleep -Seconds 1
        $retries--
    }
}

if ($retries -le 0) {
    Write-Host " Timeout waiting for Postgres." -ForegroundColor Red
    exit 1
}

if ($OnlyMigrations) {
    Write-Host "`n🎉 Migrations completed successfully!" -ForegroundColor Green
    exit 0
}

# ----------------------------------------------------
# 2. Start Application (Local)
# ----------------------------------------------------

# Backend
Write-Host "`n2. Starting Backend (Go)..." -ForegroundColor Yellow
$backendProcess = Get-Process -Name "main" -ErrorAction SilentlyContinue 
if (-not $backendProcess) {
    # Open new window for backend
    # Sets DATABASE_URL to localhost because it's running on host, not container network
    $backendCmd = "
        `$env:DATABASE_URL='postgres://specforge:specforge@localhost:5432/specforge?sslmode=disable';
        `$env:JWT_SECRET='supersecret';
        `$env:PORT='8080';
        Write-Host 'Starting SpecForge Backend...';
        cd backend;
        go run ./cmd/server
    "
    Start-Process powershell -ArgumentList "-NoExit", "-Command", "& { $backendCmd }"
    Write-Host " Backend launched in new window" -ForegroundColor Green
} else {
    Write-Host "  Backend appears to be running already" -ForegroundColor Cyan
}

# Frontend
Write-Host "`n3. Starting Frontend (Vite)..." -ForegroundColor Yellow
try {
    $frontendPortOpen = $false
    try {
        $tcp = New-Object System.Net.Sockets.TcpClient
        $tcp.Connect("localhost", 5173)
        $tcp.Close()
        $frontendPortOpen = $true
    } catch {}

    if (-not $frontendPortOpen) {
        $frontendCmd = "
            Write-Host 'Starting SpecForge Frontend...';
            cd frontend;
            npm run dev
        "
        Start-Process powershell -ArgumentList "-NoExit", "-Command", "$frontendCmd"
        Write-Host " Frontend launched in new window" -ForegroundColor Green
    } else {
        Write-Host "  Frontend appears to be running on port 5173" -ForegroundColor Cyan
    }
} catch {
    Write-Host " Failed to check frontend status" -ForegroundColor Red
}

# ----------------------------------------------------
# 3. Verification
# ----------------------------------------------------
Write-Host "`n4. Verifying Endpoints (Waiting for startup)..." -ForegroundColor Yellow

# Poll for Backend Health
$maxRetries = 20
$backendReady = $false
Write-Host -NoNewline "   Waiting for Backend..."
for ($i=0; $i -lt $maxRetries; $i++) {
    try {
        $req = Invoke-WebRequest -Uri "http://localhost:8080/health" -Method Get -UseBasicParsing -TimeoutSec 1 -ErrorAction Stop
        if ($req.StatusCode -eq 200) {
            $backendReady = $true
            break
        }
    } catch {
        Start-Sleep -Seconds 1
        Write-Host -NoNewline "."
    }
}
Write-Host ""
if ($backendReady) { Write-Host "    Backend is UP" -ForegroundColor Green } else { Write-Host "    Backend failed to start within timeout" -ForegroundColor Red }

# Poll for Frontend
$frontendReady = $false
Write-Host -NoNewline "   Waiting for Frontend..."
for ($i=0; $i -lt $maxRetries; $i++) {
    try {
        $tcp = New-Object System.Net.Sockets.TcpClient
        $tcp.Connect("localhost", 5173)
        $tcp.Close()
        $frontendReady = $true
        break
    } catch {
        Start-Sleep -Seconds 1
        Write-Host -NoNewline "."
    }
}
Write-Host ""
if ($frontendReady) { Write-Host "    Frontend is UP" -ForegroundColor Green } else { Write-Host "    Frontend failed to start within timeout" -ForegroundColor Red }

Write-Host "`n----------------------------------------"
if ($backendReady -and $frontendReady) {
    Write-Host " Development Environment is READY!" -ForegroundColor Green
    Write-Host "   Frontend: http://localhost:5173"
    Write-Host "   Backend:  http://localhost:8080"
} else {
    Write-Host "  Some services failed to start correctly." -ForegroundColor Yellow
}
Write-Host "----------------------------------------`n"
