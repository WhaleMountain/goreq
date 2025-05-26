// Command goreq provides a MCP server that fetches web content and converts it to markdown
package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/WhaleMountain/goreq/internal/handler"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// シグナルを受け取る
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var transport string
	flag.StringVar(&transport, "t", "stdio", "Transport type (stdio or sse)")
	flag.StringVar(&transport, "transport", "stdio", "Transport type (stdio or sse)")
	flag.Parse()

	// Create MCP server
	s := server.NewMCPServer(
		"GoReq",
		"1.0.0",
	)

	// Add tool
	tool := mcp.NewTool("get_url_content_for_markdown",
		mcp.WithDescription("URLのコンテンツの取得"),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("WebSite URL"),
		),
	)

	// New handler
	handler, err := handler.NewHandler()
	if err != nil {
		log.Fatalf("Failed to create handler: %v", err)
	}
	defer handler.Cleanup()

	go func() {
		<-ctx.Done()
		log.Println("Shutting down mcp server...")
		handler.Cleanup()
		os.Exit(0)
	}()

	// Add tool handler
	s.AddTool(tool, handler.HandleRequest)

	if transport == "sse" {
		sseServer := server.NewSSEServer(s, server.WithBaseURL("http://localhost:9000"))
		log.Printf("SSE Server listening on :9000")
		if err := sseServer.Start(":9000"); err != nil {
			log.Fatalf("Failed to start sse server: %v", err)
		}

	} else if transport == "http" {
		log.Printf("Starting HTTP Server on :9000")
		httpServer := server.NewStreamableHTTPServer(s)
		log.Printf("HTTP Server listening on :9000")
		if err := httpServer.Start(":9000"); err != nil {
			log.Fatalf("Failed to start HTTP server: %v\n", err)
		}

	} else if transport == "stdio" {
		log.Printf("Starting stdio Server")
		if err := server.ServeStdio(s); err != nil {
			log.Fatalf("Failed to start stdio server: %v\n", err)
		}

	} else {
		log.Fatalf("Invalid transport type: %s. Must be 'stdio' or 'sse'", transport)
	}
}
