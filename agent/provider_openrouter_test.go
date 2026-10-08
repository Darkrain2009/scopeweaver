package agent

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Autumn-27/norma/llm"
)

// A pasted OpenRouter key is sent as a Bearer token to the Chat Completions
// endpoint, and a normal text reply counts as a working connection.
func TestOpenRouterConnectionUsesPastedKey(t *testing.T) {
	const apiKey = "sk-or-pasted-key"
	var body glmWireRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/chat/completions" {
			t.Errorf("request=%s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer "+apiKey {
			t.Errorf("authorization=%q", r.Header.Get("Authorization"))
		}
		raw, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read body: %v", err)
		}
		if err := json.Unmarshal(raw, &body); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"choices":[{"message":{"role":"assistant","content":"OK"},"finish_reason":"stop"}],"usage":{"prompt_tokens":8,"completion_tokens":1}}`)
	}))
	defer srv.Close()

	cfg := Config{
		Format:          llm.FormatOpenAI,
		BaseURL:         srv.URL + "/api/v1",
		APIKey:          apiKey,
		Model:           "z-ai/glm-5.3-flash",
		ThinkingType:    "enabled",
		ReasoningEffort: "max",
		Stream:          false,
		Retry:           RetryConfig{ConnectAttempts: -1, EmptyAttempts: -1},
	}
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	_, reply, err := TestConnection(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(reply, "OK") {
		t.Fatalf("reply=%q", reply)
	}
	if body.Model != "z-ai/glm-5.3-flash" || body.Thinking.Type != "enabled" || body.ReasoningEffort != "max" {
		t.Fatalf("request=%+v", body)
	}
}
