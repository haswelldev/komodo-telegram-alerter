package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var client = &http.Client{Timeout: 10 * time.Second}

type sendPayload struct {
	ChatID                string `json:"chat_id"`
	Text                  string `json:"text"`
	ParseMode             string `json:"parse_mode"`
	DisableWebPagePreview bool   `json:"disable_web_page_preview"`
}

type apiResponse struct {
	OK          bool   `json:"ok"`
	Description string `json:"description"`
}

// Send posts an HTML-formatted message to a Telegram chat.
func Send(token, chatID, html string) error {
	payload := sendPayload{
		ChatID:                chatID,
		Text:                  html,
		ParseMode:             "HTML",
		DisableWebPagePreview: true,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
	resp, err := client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("http post: %w", err)
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	var apiResp apiResponse
	if err := json.Unmarshal(raw, &apiResp); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	if !apiResp.OK {
		log.Printf("[ERROR] Telegram API error: %s\n--- message ---\n%s", apiResp.Description, html)
		return fmt.Errorf("telegram API: %s", apiResp.Description)
	}
	return nil
}
