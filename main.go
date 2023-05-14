package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Poul-george/go_chatgpt_api/config"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Client struct {
	http.Client
}

type Request struct {
	Model       string `json:"model"`
	Prompt      string `json:"prompt"`
	Temperature string `json:"temperature"`
	MaxTokens   string `json:"max_tokens"`
}

type CompletionResponse struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

func main() {
	fmt.Println("質問したい内容を書いてください")
	scanner := bufio.NewScanner(os.Stdin)

	// default question text
	text := "nftは今後どれくらい伸びるのでしょうか？"
	if scanner.Scan() {
		text = scanner.Text()
	}

	req, err := NewRequest(text)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(req)

	c := Client{}
	resp, err := c.Do(req)
	if err != nil {
		fmt.Println(resp.Status)
		log.Fatal(err)
		return
	}
	fmt.Println(resp.Status)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return
	}

	var completionResponse CompletionResponse

	if err := json.Unmarshal(body, &completionResponse); err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println()
	fmt.Println(completionResponse.Choices)
	fmt.Println()
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewRequest(text string) (*http.Request, error) {
	apiKey := config.GetConfig().APIKey
	fmt.Println(apiKey)

	requestBodyBytes, err := json.Marshal(RequestBody(text))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

func RequestBody(t string) Request {
	return Request{
		Model:       "text-davinci-003",
		Prompt:      t,
		Temperature: "1",
		MaxTokens:   "100",
	}
}
