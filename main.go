package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/hamisionesmus/project4/cache"
	"github.com/hamisionesmus/project4/server"
)

func main() {
	// Parse command line flags
	config := parseFlags()

	// Initialize logger
	logger := log.New(os.Stdout, "[CACHE] ", log.LstdFlags)

	// Create cache instance
	cacheInstance := cache.NewCache(cache.Config{
		MaxMemory:     config.MaxMemory,
		EvictionPolicy: config.EvictionPolicy,
	})

	// Start cache cleanup routine
	go cacheInstance.StartCleanup()

	// Create TCP server
	tcpServer := server.NewTCPServer(cacheInstance, logger)

	// Start TCP server
	go func() {
		logger.Printf("Starting TCP server on %s:%d", config.Host, config.Port)
		if err := tcpServer.Start(fmt.Sprintf("%s:%d", config.Host, config.Port)); err != nil {
			logger.Fatalf("TCP server failed: %v", err)
		}
	}()

	// Start HTTP monitoring server if enabled
	if config.HTTPPort > 0 {
		httpServer := server.NewHTTPServer(cacheInstance, logger)
		go func() {
			logger.Printf("Starting HTTP server on %s:%d", config.Host, config.HTTPPort)
			if err := httpServer.Start(fmt.Sprintf("%s:%d", config.Host, config.HTTPPort)); err != nil {
				logger.Fatalf("HTTP server failed: %v", err)
			}
		}()
	}

	// Wait for interrupt signal
	waitForShutdown()

	// Graceful shutdown
	logger.Println("Shutting down servers...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		tcpServer.Shutdown(ctx)
	}()

	wg.Wait()
	logger.Println("Servers shut down gracefully")
}

type Config struct {
	Host           string
	Port           int
	HTTPPort       int
	MaxMemory      string
	EvictionPolicy string
}

func parseFlags() *Config {
	// Simple flag parsing (in real implementation, use flag package)
	host := getEnv("CACHE_HOST", "0.0.0.0")
	port := getEnvInt("CACHE_PORT", 8080)
	httpPort := getEnvInt("CACHE_HTTP_PORT", 8081)
	maxMemory := getEnv("CACHE_MAX_MEMORY", "1GB")
	evictionPolicy := getEnv("CACHE_EVICTION_POLICY", "lru")

	return &Config{
		Host:           host,
		Port:           port,
		HTTPPort:       httpPort,
		MaxMemory:      maxMemory,
		EvictionPolicy: evictionPolicy,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func waitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}

// Placeholder implementations (would be in separate files in real project)

type Cache struct {
	data map[string]interface{}
	mu   sync.RWMutex
}

func NewCache(config cache.Config) *Cache {
	return &Cache{
		data: make(map[string]interface{}),
	}
}

func (c *Cache) StartCleanup() {
	// Implementation for cleanup routine
}

type TCPServer struct{}

func NewTCPServer(cache *Cache, logger *log.Logger) *TCPServer {
	return &TCPServer{}
}

func (s *TCPServer) Start(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go s.handleConnection(conn)
	}
}

func (s *TCPServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	// Handle Redis protocol commands
	conn.Write([]byte("+OK\r\n"))
}

func (s *TCPServer) Shutdown(ctx context.Context) error {
	return nil
}

type HTTPServer struct{}

func NewHTTPServer(cache *Cache, logger *log.Logger) *HTTPServer {
	return &HTTPServer{}
}

func (s *HTTPServer) Start(addr string) error {
	// HTTP server implementation would go here
	return nil
}