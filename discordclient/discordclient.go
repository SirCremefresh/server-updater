package discordclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Severity int

const (
	Info Severity = iota
	Error
)

// AptClient is cool
type DiscordClient struct {
	botName         string
	errorWebhookURL string
	infoWebhookURL  string
}

// New is cool
func New(botName, errorWebhookURL, infoWebhookURL string) *DiscordClient {
	return &DiscordClient{
		botName:         botName,
		errorWebhookURL: errorWebhookURL,
		infoWebhookURL:  infoWebhookURL,
	}
}

func (client *DiscordClient) LogFatalIgnoreError(message string) {
	//	err := client.Send(message, Error)
	//	if err != nil {
	//		log.Printf("error while sending discord message. err: %v\n", err)
	//	}
	log.Fatalf("error: %s", message)
}

func (client *DiscordClient) LogInfoIgnoreError(message string) {
	log.Printf("info: %s\n", message)
	//	err := client.Send(message, Info)
	//	if err != nil {
	//		log.Printf("error while sending discord message. err: %v\n", err)
	//	}
}

func (client *DiscordClient) Send(message string, severity Severity) error {
	var webhookURL string
	if severity == Info {
		webhookURL = client.infoWebhookURL
	} else {
		webhookURL = client.errorWebhookURL
	}

	requestBody, err := json.Marshal(map[string]string{
		"username": client.botName,
		"content":  message,
	})
	if err != nil {
		return fmt.Errorf("could not json marshal webhook request body. err: %v", err)
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("could not make webhook request. err: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read webhook response body. err: %v", err)
	}
	if resp.StatusCode != 204 {
		return fmt.Errorf("request to webhook failed. status: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}
