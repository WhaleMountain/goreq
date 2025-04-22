package handler

import (
	"context"
	"fmt"
	"goreq/browser"
	gourl "net/url"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/mark3labs/mcp-go/mcp"
)

type Handler struct {
	browser *browser.Browser
}

func NewHandler() (*Handler, error) {
	browser, err := browser.NewBrowser()
	if err != nil {
		return nil, err
	}
	handler := &Handler{browser: browser}
	return handler, nil
}

func (h *Handler) HandleRequest(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	url, ok := request.Params.Arguments["url"].(string)
	if !ok {
		return nil, fmt.Errorf("url must be a string")
	}
	if _, err := gourl.ParseRequestURI(url); err != nil {
		return nil, fmt.Errorf("invalid URL: %s", url)
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

func (h *Handler) Cleanup() {
	if h.browser != nil {
		h.browser.Close()
	}
}