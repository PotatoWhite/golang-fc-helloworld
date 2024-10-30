
## 프로젝트 구조

```
golang-fc-helloworld/
├── cmd/
│   └── golang-fc-helloworld/
│       └── main.go       # 메인 프로그램 진입점
├── pkg/
│   └── function_call/
│       └── chat.go       # ChatGPT API와의 통신 처리 로직
├── Makefile               # 빌드 및 실행을 위한 Makefile
└── go.mod                 # Go 모듈 설정 파일
```

---

## 1단계: `go.mod` 파일 초기화

```bash
go mod init golang-fc-helloworld
go mod tidy
```

---

## 2단계: `pkg/function_call/chat.go` 작성

이 파일에서는 ChatGPT API와 통신하는 로직을 구현합니다.

```go
package functioncall

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// ChatGPT 요청 구조체
type FunctionCallRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content,omitempty"`
}

// ChatGPT 응답 구조체
type ChatGPTResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

// Chat 함수: ChatGPT와 통신
func Chat(userMessage string) (string, error) {
	// 요청 본문 생성
	requestBody := FunctionCallRequest{
		Model: "gpt-4-0613",
		Messages: []Message{
			{Role: "system", Content: "You are a helpful assistant."},
			{Role: "user", Content: userMessage},
		},
	}

	// JSON 직렬화
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	// HTTP 요청 생성
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// 헤더 설정
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))

	// HTTP 클라이언트 생성 및 요청 수행
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 응답 본문 읽기
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// 응답 데이터 파싱
	var chatResponse ChatGPTResponse
	if err := json.Unmarshal(body, &chatResponse); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// 응답 반환
	if len(chatResponse.Choices) > 0 {
		return chatResponse.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no response from ChatGPT")
}
```

---

## 3단계: `cmd/golang-fc-helloworld/main.go` 작성

이 파일에서는 HTTP 서버를 구현하고, 사용자가 요청을 보낼 수 있도록 합니다.

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"golang-fc-helloworld/pkg/function_call"
	"io"
)

// /chat 엔드포인트 핸들러
func chatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	userMessage := string(body)
	reply, err := functioncall.Chat(userMessage)
	if err != nil {
		http.Error(w, "Failed to get response from ChatGPT: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"reply": "%s"}`, reply)))
}

func main() {
	http.HandleFunc("/chat", chatHandler)

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

---

## 4단계: Makefile 작성

아래는 프로젝트의 빌드 및 실행을 자동화하는 Makefile입니다.

```Makefile
.PHONY: all build run clean

APP_NAME := golang-fc-helloworld
BUILD_DIR := bin

all: build

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/$(APP_NAME)

run: build
	$(BUILD_DIR)/$(APP_NAME)

clean:
	rm -rf $(BUILD_DIR)
```

---

## 5단계: 실행 방법

1. **API 키 설정**
   ```bash
   export OPENAI_API_KEY="your-openai-api-key"
   ```

2. **빌드 및 실행**
   ```bash
   make run
   ```

3. **테스트**
   서버가 실행된 후, 다른 터미널에서 다음 명령어를 실행합니다.

   ```bash
   curl -X POST http://localhost:8080/chat -d "Hello, how are you?"
   ```

---

## 6단계: 테스트 결과

정상적으로 동작하면 아래와 같은 응답을 받을 수 있습니다.

```json
{
  "reply": "I'm good, thank you! How can I assist you today?"
}
```

---

## 7단계: 정리 및 확장

- **기능 확장**: 요청에 따라 더 다양한 기능을 구현해 보세요.
- **에러 처리 개선**: 네트워크 오류나 API 호출 실패에 대한 로직을 추가하세요.
- **프론트엔드 연동**: 간단한 HTML 페이지를 만들어 이 서버와 연동해 보세요.

---

## 결론

이번 실습에서는 Go 언어의 표준 라이브러리와 API 통신을 사용해 간단한 HTTP 서버를 구축하고 ChatGPT API와 연동하는 방법을 배웠습니다. 이 구조를 바탕으로 더 복잡한 AI 애플리케이션을 개발할 수 있습니다. 🚀