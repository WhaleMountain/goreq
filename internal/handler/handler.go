// Package handler provides handlers for MCP tools
package handler

import (
	"context"
	"fmt"
	gourl "net/url"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/WhaleMountain/goreq/internal/browser"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Handler handles MCP requests
type Handler struct {
	browser *browser.Browser
}

// NewHandler creates a new handler
func NewHandler() (*Handler, error) {
	browser, err := browser.NewBrowser()
	if err != nil {
		return nil, err
	}
	handler := &Handler{browser: browser}
	return handler, nil
}

type ToolArgs struct {
	URL string `json:"url" jsonschema:"required,description=WebSite URL"`
}

func (h *Handler) HandleRequest(ctx context.Context, request *mcp.CallToolRequest, args ToolArgs) (*mcp.CallToolResult, any, error) {
	if _, err := gourl.ParseRequestURI(args.URL); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("invalid URL: %s", args.URL)},
			},
			IsError: true,
		}, nil, nil
	}

	content, err := h.browser.GetContent(args.URL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get content: %v", err)
	}

	markdown, err := htmltomarkdown.ConvertString(content)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert markdown: %v", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: markdown},
		},
	}, nil, nil
}

// Cleanup cleans up resources
func (h *Handler) Cleanup() {
	if h.browser != nil {
		h.browser.Close()
	}
}
