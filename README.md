# WebSocket Stress Test

A Go application for running websocket stress tests, showing real-time connection statistics via a web interface.

## Directory Structure

```
websocket-stress/
├── cmd/
│   └── websocket-stress/      # Application entry point
│       └── main.go            # Main application code
├── internal/
│   └── server/                # Server implementation
│       └── server.go          # HTTP and WebSocket server handlers
├── pkg/
│   └── websocket/             # Reusable websocket components
│       └── client.go          # WebSocket client implementation
├── web/
│   ├── static/                # Static assets (CSS, JS, images)
│   │   └── style.css          # Additional CSS styles
│   └── templates/             # HTML templates
│       └── index.html         # Main application page
├── go.mod                     # Go module definition
├── go.sum                     # Go module checksum
└── README.md                  # This file
```

## Running the Application

To run the application, navigate to the project root and use one of these commands:

```bash
# Run directly
go run cmd/websocket-stress/main.go

# Or build and run
cd cmd/websocket-stress
go build
./websocket-stress
```

The application will be available at http://localhost:8088 from both inside and outside the container.

**Note for VS Code Remote Development users:** If you're using the devcontainer setup, you'll need to rebuild the container after changing the port configurations by using the "Rebuild Container" command in VS Code.

## Features

- Real-time websocket connection monitoring
- Connection statistics (uptime, message count, client count)
- Auto-reconnect on connection loss
- Clean, responsive UI

## Dependencies

- github.com/gorilla/websocket v1.5.0 - WebSocket implementation for Go

## For Computer Lab Testing

1. Run the server on a centralized machine that is accessible to all computers in the lab.
2. Make sure the server IP is reachable and the port 8088 is open.
3. Have all computers in the lab navigate to `http://SERVER_IP:8088` where `SERVER_IP` is the IP address of the machine running the server.
4. Each connection will automatically be tracked and displayed on all clients.

## Testing Connection Failures

You can test connection failures in several ways:

1. Disconnect the client computer from the network temporarily
2. Shut down the server while clients are connected
3. Introduce network latency or packet loss using network tools

When a connection fails, the client interface will show a red "Disconnected" status and attempt to reconnect automatically. 