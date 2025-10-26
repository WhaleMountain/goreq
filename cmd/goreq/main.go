// Command goreq provides a MCP server that fetches web content and converts it to markdown
package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/WhaleMountain/goreq/internal/handler"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	var transport string
	flag.StringVar(&transport, "t", "stdio", "Transport type (stdio or http)")
	flag.StringVar(&transport, "transport", "stdio", "Transport type (stdio or http)")
	var httpAddr string
	flag.StringVar(&httpAddr, "http", ":9000", "HTTP address for streamable HTTP transport")
	flag.Parse()

	// New handler
	h, err := handler.NewHandler()
	if err != nil {
		log.Fatalf("Failed to create handler: %v", err)
	}
	defer h.Cleanup()

	// Create MCP server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "GoReq",
		Version: "1.2.0",
	}, nil)

	// Add tool
	tool := &mcp.Tool{
		Name:        "get_url_content_for_markdown",
		Description: "URLのコンテンツの取得",
	}

	// Add tool handler
	mcp.AddTool(server, tool, h.HandleRequest)

	switch transport {
	case "http":
		httpHandler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
			return server
		}, nil)
		log.Printf("HTTP Server listening on %s", httpAddr)
		if err := http.ListenAndServe(httpAddr, httpHandler); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	case "stdio":
		log.Printf("Starting stdio Server")
		t := &mcp.StdioTransport{}
		if err := server.Run(context.Background(), t); err != nil {
			log.Fatalf("Failed to start stdio server: %v\n", err)
		}
	default:
		log.Fatalf("Invalid transport type: %s. Must be 'stdio' or 'http'", transport)
	}
}
