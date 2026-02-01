# Go Port Scanner

A simple port scanner written in Go.
 
# Features

Concurrent Scanning: Utilizes Go's goroutines for parallel port scanning
Two Scanning Modes:
Unlimited Mode: Uses semaphore-based concurrency control with automatic file descriptor limits
Limited Mode: Fixed concurrency limit of 1000 concurrent connections
Automatic Resource Management: Detects system file descriptor limits using ulimit
Timeout Handling: Configurable connection timeouts
Retry Logic: Automatically retries when hitting "too many open files" errors
Prerequisites

Go 1.16 or higher
Unix-like operating system (Linux, macOS) for ulimit support
Dependencies

go get golang.org/x/sync/semaphore
go get github.com/pieterclaerhout/go-waitgroup

# Installation

Clone or download the code
git clone <repository-url>
cd <project-directory>

Install dependencies
go mod download

Build the application
go build -o portscanner

#Usage

Basic Usage

The default configuration scans localhost (127.0.0.1) for all ports (0-65535):

./portscanner

Customization

Modify the main() function to customize scanning parameters:

ps := &PortScanner {
    host: "192.168.1.1",  // Target host
    lock: semaphore.NewWeighted(ulimit()),
}

# Key Functions

ulimit(): Retrieves system file descriptor limit
scanPort(): Attempts TCP connection to a specific port
StartUnlimited(): Scans with semaphore-based concurrency (respects system limits)
StartLimited(): Scans with fixed concurrency limit (1000 goroutines)
