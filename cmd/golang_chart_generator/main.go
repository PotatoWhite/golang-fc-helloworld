package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"golang-fc-helloworld/pkg/function_call"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	openaiAPIURL = "https://api.openai.com/v1/chat/completions"
	model        = "gpt-4-0613"
	maxTurns     = 10
)

func main() {
	// Initialize messages and functions
	messages := function_call.InitialMessage()
	functions := function_call.GetFunctionSpec()

	turns := 0
	reader := bufio.NewReader(os.Stdin)

	for turns < maxTurns {

		if turns == 0 {
			fmt.Println("안녕하세요, 저는 담당 의사입니다. 먼저 증상과 성별, 나이를 알려주시겠어요? (예: '저는 두통이 있어요. 여성, 30대')")
		}
		fmt.Print("사용자: ")
		userInput, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading user input:", err)
			continue
		}
		userInput = userInput[:len(userInput)-1] // 개행 문자 제거

		messages = append(messages, function_call.Message{
			Role:    "user",
			Content: userInput,
		})

		response, err := sendRequestToOpenAI(messages, functions)
		if err != nil {
			log.Println("Error sending request to OpenAI:", err)
			continue
		}

		processResponse(response, &messages)

		turns++
	}

	log.Println("대화를 종료합니다.")
}

func sendRequestToOpenAI(messages []function_call.Message, functions []function_call.FunctionSpec) (*function_call.ResponseBody, error) {
	requestBody := function_call.RequestBody{
		Model:               model,
		Messages:            messages,
		Functions:           functions,
		FunctionCallSetting: "auto",
	}

	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request body: %v", err)
	}

	req, err := http.NewRequest("POST", openaiAPIURL, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPEN_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var responseBody function_call.ResponseBody
	if err := json.Unmarshal(body, &responseBody); err != nil {
		return nil, fmt.Errorf("error unmarshalling response body: %v", err)
	}

	return &responseBody, nil
}

func processResponse(response *function_call.ResponseBody, messages *[]function_call.Message) {
	if len(response.Choices) == 0 {
		log.Println("No response from the API.")
		return
	}

	responseContent := response.Choices[0].Message.Content
	functionCall := response.Choices[0].Message.FunctionCall

	fmt.Println("의사의 응답:", responseContent)

	if functionCall.Name != "" {
		function_call.HandleFunctionCall(functionCall)
	}

	*messages = append(*messages, function_call.Message{
		Role:    "assistant",
		Content: responseContent,
	})
}
