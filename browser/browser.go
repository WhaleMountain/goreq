package browser

import (
	"fmt"

	"github.com/playwright-community/playwright-go"
)

type Browser struct {
	pw      *playwright.Playwright
	browser playwright.Browser
}

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

func (b *Browser) Close() {
	if b.browser != nil {
		b.browser.Close()
	}
	if b.pw != nil {
		b.pw.Stop()
	}
}

func (b *Browser) GetContent(url string) (string, error) {
	page, err := b.browser.NewPage(playwright.BrowserNewPageOptions{
		UserAgent: playwright.String("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create page: %v", err)
	}
	defer page.Close()

	if _, err = page.Goto(url, playwright.PageGotoOptions{WaitUntil: playwright.WaitUntilStateLoad}); err != nil {
		return "", fmt.Errorf("could not goto: %v", err)
	}

	content, err := page.Content()
	if err != nil {
		return "", err
	}

	return content, nil
}
