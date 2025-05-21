// Package handler provides handlers for MCP tools
package handler

import (
	"context"
	"fmt"
	gourl "net/url"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/WhaleMountain/goreq/internal/browser"
	"github.com/mark3labs/mcp-go/mcp"
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

// HandleRequest handles MCP tool requests
func (h *Handler) HandleRequest(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	url, err := request.RequireString("url")
	if err != nil {
		return mcp.NewToolResultError((err.Error())), nil
	}
	if _, err := gourl.ParseRequestURI(url); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("invalid URL: %s", url)), nil
	}

	content, err := h.browser.GetContent(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get content: %v", err)
	}

	markdown, err := htmltomarkdown.ConvertString(content)
	if err != nil {
		return nil, fmt.Errorf("failed to convert markdown: %v", err)
	}

	return mcp.NewToolResultText(markdown), nil
}

// Cleanup cleans up resources
func (h *Handler) Cleanup() {
	if h.browser != nil {
		h.browser.Close()
	}
}
