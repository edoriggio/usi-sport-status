package main

import (
	"os"
	"fmt"
	"bytes"
	"errors"
	"strings"
	"net/http"
	"encoding/json"

	"github.com/joho/godotenv"
)

type webhookReqBody struct {
	Message struct {
		Text string `json:"text"`
		Chat struct {
			ID int64  `json:"id"`
		} `json:"chat"`
	} `json:"message"`
}

type sendMessageReqBody struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

func Handler(res http.ResponseWriter, req *http.Request) {
	body := &webhookReqBody{}

	if err := json.NewDecoder(req.Body).Decode(body); err != nil {
		fmt.Println("could not decode request body", err)
		return
	}

	if !strings.Contains(strings.ToLower(body.Message.Text), "status") {
		return
	}

	if err := sendMessage(body.Message.Chat.ID); err != nil {
		fmt.Println("error in sending reply:", err)
		return
	}

	fmt.Println("reply sent")
}

func sendMessage(chatID int64) error {
	isAvailable := checkIfExpired()

	message := "SPOTS AVAILABLE!!!"

	if !isAvailable {
		message = "No spots available"
	}

	reqBody := &sendMessageReqBody{
		ChatID: chatID,
		Text:   message,
	}

	reqBytes, err := json.Marshal(reqBody)

	if err != nil {
		return err
	}

	godotenv.Load("../.env")

	token := os.Getenv("TELEGRAM_BOT")
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	res, err := http.Post(url, "application/json", bytes.NewBuffer(reqBytes))

	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return errors.New("unexpected status" + res.Status)
	}

	return nil
}

func main() {
	http.ListenAndServe(":3000", http.HandlerFunc(Handler))
}
