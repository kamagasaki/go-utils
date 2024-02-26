package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

func SendMessage(apiKey string, chatID interface{}, message string) error {
	// Check the type of chatID and convert it accordingly
	chatIDValue := ""
	switch v := chatID.(type) {
	case int64:
		chatIDValue = strconv.FormatInt(v, 10)
	case string:
		chatIDValue = v
	default:
		return fmt.Errorf("unsupported chatID type: %T", chatID)
	}

	requestBody, err := json.Marshal(map[string]interface{}{
		"chat_id":    chatIDValue,
		"text":       message,
		"parse_mode": "Markdown",
	})
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.telegram.org/bot"+apiKey+"/sendMessage", bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func SendDocument(token, chatID interface{}, filePath, caption string) error {
	// Check the type of chatID and convert it accordingly
	chatIDValue := ""
	switch v := chatID.(type) {
	case int64:
		chatIDValue = strconv.FormatInt(v, 10)
	case string:
		chatIDValue = v
	default:
		return fmt.Errorf("unsupported chatID type: %T", chatID)
	}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a new multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the file field
	part, err := writer.CreateFormFile("document", filePath)
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	// Add other form fields
	_ = writer.WriteField("chat_id", chatIDValue)
	_ = writer.WriteField("caption", caption)
	_ = writer.WriteField("parse_mode", "Markdown")

	// Close the multipart writer
	err = writer.Close()
	if err != nil {
		return err
	}

	// Create the request
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendDocument", token)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
