# Local Network Monitor

A cross-platform desktop application that monitors and displays open TCP ports on your local machine, built with Go (backend) and Flutter (frontend).

## Features

- **Real-time Port Scanning**: Displays all open TCP ports on your system
- **Process Detection**: Shows which process is using each port
- **Cross-Platform**: Works on both Linux and Windows
- **Modern UI**: Beautiful Material 3 design with dark theme
- **Auto-Launch Backend**: Frontend automatically starts the Go backend server

## Architecture

- **Backend**: Go HTTP server (port 8080) that scans system ports
  - Linux: Uses `/proc/net/tcp` for port scanning
  - Windows: Uses `netstat` and `tasklist` commands
- **Frontend**: Flutter desktop application with modern UI
- **Communication**: REST API over HTTP

## Prerequisites

- **Go** 1.21 or higher
- **Flutter** 3.0 or higher
- **Linux** or **Windows** operating system

## Project Structure

```
localnetworkmonitor/
├── backend/
│   ├── main.go           # HTTP server and API endpoints
│   ├── scanner.go        # Port scanning logic (Linux & Windows)
│   └── go.mod           # Go module definition
├── frontend/
│   └── local_network_monitor/
│       ├── lib/
│       │   └── main.dart    # Flutter application
│       └── pubspec.yaml     # Flutter dependencies
├── build/               # Build output directory
├── run.sh              # Linux/macOS build & run script
├── run.bat             # Windows build & run script
└── README.md
```

## Installation & Running

### Linux/macOS

1. Clone the repository
2. Make the run script executable:
   ```bash
   chmod +x run.sh
   ```
3. Run the application:
   ```bash
   ./run.sh
   ```

### Windows

1. Clone the repository
2. Run the application:
   ```cmd
   run.bat
   ```

The scripts will:
1. Build the Go backend to `build/network_monitor_backend` (or `.exe` on Windows)
2. Build the Flutter application in release mode
3. Copy all files to the `build/` directory

## Manual Build

### Backend Only

```bash
cd backend
go build -o ../build/network_monitor_backend
```

### Frontend Only

```bash
cd frontend/local_network_monitor
flutter build linux --release  # or 'windows' on Windows
```

## API Endpoints

### GET /ports

Returns a JSON array of open ports with their associated processes.

**Response Example:**
```json
[
  {
    "port": 80,
    "process": "nginx"
  },
  {
    "port": 5432,
    "process": "postgres"
  },
  {
    "port": 8080,
    "process": "network_monitor_backend"
  }
]
```

## How It Works

1. **Backend Launch**: When the Flutter app starts, it automatically launches the Go backend binary
2. **Backend Ready Check**: The app waits for the backend to be ready (max 6 seconds with 200ms intervals)
3. **Port Scanning**: 
   - On Linux: Reads `/proc/net/tcp` and maps inodes to processes via `/proc/[pid]/fd`
   - On Windows: Executes `netstat -ano` and `tasklist` to get port and process information
4. **Data Display**: Frontend fetches data from the backend API and displays it in a modern card-based UI

## UI Features

- **Port Count Badge**: Shows total number of open ports
- **Refresh Button**: Manually refresh the port list
- **Loading State**: Animated spinner while scanning
- **Error Handling**: Displays errors if backend is unreachable
- **Empty State**: Shows message when no ports are found
- **Active Status**: Each port displays an "ACTIVE" badge

## Permissions

### Linux
- May require elevated privileges to access certain process information
- Run with `sudo` if you need complete process details

### Windows
- Run as Administrator for complete process information
- Standard user can see ports but may have limited process details

## Development

### Backend Development

```bash
cd backend
go run .
```

The server will start on `http://localhost:8080`

### Frontend Development

```bash
cd frontend/local_network_monitor
flutter run -d linux  # or -d windows
```

## Troubleshooting

**Backend not starting:**
- Ensure port 8080 is not already in use
- Check if the backend binary exists in `build/` directory
- Verify Go is installed: `go version`

**Frontend build errors:**
- Run `flutter doctor` to check Flutter installation
- Run `flutter pub get` to install dependencies

**No ports showing:**
- On Linux: Check if `/proc/net/tcp` is accessible
- On Windows: Ensure `netstat` and `tasklist` commands work
- Try running with elevated privileges

## License

This project is provided as-is for educational and personal use.

## Technologies Used

- **Go**: Backend server and system port scanning
- **Flutter**: Cross-platform desktop UI framework
- **Material 3**: Modern design system
- **Riverpod**: State management for Flutter
