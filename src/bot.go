// Copyright 2021 Edoardo Riggio
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

	if err := sendMessage(body.Message.Chat.ID, body.Message.Text); err != nil {
		fmt.Println("error in sending reply:", err)
		return
	}

	fmt.Println("reply sent")
}

func sendMessage(chatID int64, body string) error {
	isAvailable := checkIfExpired(body)

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
