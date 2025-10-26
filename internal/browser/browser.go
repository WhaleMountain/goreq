// Package browser provides browser functionality for accessing web content
package browser

import (
	"fmt"
	"sync"

	"github.com/playwright-community/playwright-go"
)

// Browser represents a browser instance
type Browser struct {
	pw      *playwright.Playwright
	browser playwright.Browser
	mu      sync.Mutex // 追加: スレッドセーフ用
}

// NewBrowser creates a new browser instance
func NewBrowser() (*Browser, error) {
	pw, err := playwright.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to create Playwright: %v", err)
	}
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to launch browser: %v", err)
	}
	return &Browser{
		pw:      pw,
		browser: browser,
	}, nil
}

// Close closes the browser instance
func (b *Browser) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.browser != nil {
		if err := b.browser.Close(); err != nil {
			fmt.Printf("failed to close browser: %v\n", err)
		}
		b.browser = nil
	}
	if b.pw != nil {
		if err := b.pw.Stop(); err != nil {
			fmt.Printf("failed to stop Playwright: %v\n", err)
		}
		b.pw = nil
	}
}

// GetContent fetches the HTML content from the specified URL
func (b *Browser) GetContent(url string) (string, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	page, err := b.browser.NewPage(playwright.BrowserNewPageOptions{
		UserAgent: playwright.String("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create page: %v", err)
	}
	defer func() {
		if err := page.Close(); err != nil {
			fmt.Printf("failed to close page: %v\n", err)
		}
	}()

	if _, err = page.Goto(url, playwright.PageGotoOptions{WaitUntil: playwright.WaitUntilStateLoad}); err != nil {
		return "", fmt.Errorf("could not goto: %v", err)
	}

	content, err := page.Content()
	if err != nil {
		return "", err
	}

	return content, nil
}
