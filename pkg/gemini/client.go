package gemini

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"agentic/pkg/contract"
	"agentic/pkg/jev"
)

const (
	DefaultBaseURL = "https://generativelanguage.googleapis.com/v1beta"
	DefaultModel   = "gemini-2.0-flash"
)

// Compile-time assertion: *Client satisfies the jev.Client interface.
var _ jev.Client = (*Client)(nil)

// Client interacts with Google Cloud Gemini API
type Client struct {
	apiKey     string
	model      string
	baseURL    string
	httpClient *http.Client
}

func NewClient(apiKey, model string) *Client {
	if model == "" {
		model = DefaultModel
	}
	return &Client{
		apiKey:  apiKey,
		model:   model,
		baseURL: DefaultBaseURL,
		httpClient: &http.Client{
			Timeout: 45 * time.Second,
		},
	}
}

// System 2: GenerateCode / Plan (Creative synthesis)
func (c *Client) GenerateContent(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	if c.apiKey == "" {
		return "", fmt.Errorf("gemini API key is not configured")
	}

	reqBody := map[string]any{
		"contents": []map[string]any{
			{
				"role": "user",
				"parts": []map[string]any{
					{"text": userPrompt},
				},
			},
		},
	}

	if systemPrompt != "" {
		reqBody["systemInstruction"] = map[string]any{
			"parts": []map[string]any{
				{"text": systemPrompt},
			},
		}
	}

	respBytes, err := c.post(ctx, c.model, reqBody)
	if err != nil {
		return "", fmt.Errorf("gemini generateContent error: %w", err)
	}

	return extractText(respBytes)
}

// System 1: DecideRoute using Gemini Structured Outputs (JSON mode)
func (c *Client) DecideRoute(ctx context.Context, tc *contract.TaskContract, evidence []contract.Evidence, iteration, maxIterations int) (*jev.RoutingDecision, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("gemini API key is not configured")
	}

	prompt := fmt.Sprintf(`Task Objective: %s
Iteration: %d / %d
Verifiable Evidence Collected:
`, tc.Objective, iteration, maxIterations)

	for _, ev := range evidence {
		status := "PASS"
		if !ev.Passed {
			status = "FAIL"
		}
		prompt += fmt.Sprintf("- [%s] (%s): %s -> %s\n", status, ev.Type, ev.Claim, ev.Details)
	}

	systemPrompt := `You are an ultra-fast System 1 AI router.
Analyze the evidence against the objective and return a JSON decision adhering strictly to this schema:
{
  "action": "PROCEED_TO_GATE" | "RETRY_REPAIR" | "HALT_EXHAUSTED",
  "confidence": <float between 0.0 and 1.0>,
  "reason": "<concise explanation>",
  "suggested_fixes": ["<fix1>", "<fix2>"]
}
Rules:
- If all evidence PASS, action MUST be "PROCEED_TO_GATE".
- If any evidence FAIL and iteration < maxIterations, action MUST be "RETRY_REPAIR" with suggested fixes.
- If iteration >= maxIterations, action MUST be "HALT_EXHAUSTED".
`

	reqBody := map[string]any{
		"systemInstruction": map[string]any{
			"parts": []map[string]any{{"text": systemPrompt}},
		},
		"contents": []map[string]any{
			{
				"parts": []map[string]any{{"text": prompt}},
			},
		},
		"generationConfig": map[string]any{
			"temperature":      0.1,
			"responseMimeType": "application/json",
		},
	}

	respBytes, err := c.post(ctx, c.model, reqBody)
	if err != nil {
		return nil, fmt.Errorf("gemini route decision failed: %w", err)
	}

	rawText, err := extractText(respBytes)
	if err != nil {
		return nil, err
	}

	var decision jev.RoutingDecision
	if err := json.Unmarshal([]byte(rawText), &decision); err != nil {
		return nil, fmt.Errorf("failed to parse structured gemini decision JSON: %w (raw: %s)", err, rawText)
	}

	return &decision, nil
}

// System 1: ScoreVerification using Gemini Structured Outputs
func (c *Client) ScoreVerification(ctx context.Context, tc *contract.TaskContract, candidateDiff string, evidence []contract.Evidence) (*jev.VerifierScore, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("gemini API key is not configured")
	}

	evidenceSummary := ""
	for _, ev := range evidence {
		status := "PASS"
		if !ev.Passed {
			status = "FAIL"
		}
		evidenceSummary += fmt.Sprintf("- [%s] (%s): %s -> %s\n", status, ev.Type, ev.Claim, ev.Details)
	}

	systemPrompt := `You are an independent adversarial verifier checking if the candidate has disproof.
Return a JSON score strictly adhering to this schema:
{
  "passed": <boolean>,
  "confidence": <float between 0.0 and 1.0>,
  "disproof_found": <boolean>,
  "rejection_reason": "<optional reason if disproof found>"
}
If any deterministic check failed, disproof_found MUST be true and passed MUST be false.
`

	prompt := fmt.Sprintf("Objective: %s\nCandidate Diff:\n%s\nEvidence:\n%s", tc.Objective, candidateDiff, evidenceSummary)

	reqBody := map[string]any{
		"systemInstruction": map[string]any{
			"parts": []map[string]any{{"text": systemPrompt}},
		},
		"contents": []map[string]any{
			{"parts": []map[string]any{{"text": prompt}}},
		},
		"generationConfig": map[string]any{
			"temperature":      0.1,
			"responseMimeType": "application/json",
		},
	}

	respBytes, err := c.post(ctx, c.model, reqBody)
	if err != nil {
		return nil, err
	}

	rawText, err := extractText(respBytes)
	if err != nil {
		return nil, err
	}

	var score jev.VerifierScore
	if err := json.Unmarshal([]byte(rawText), &score); err != nil {
		return nil, fmt.Errorf("failed to parse gemini verifier score JSON: %w (raw: %s)", err, rawText)
	}

	return &score, nil
}

func (c *Client) post(ctx context.Context, model string, body map[string]any) ([]byte, error) {
	url := fmt.Sprintf("%s/models/%s:generateContent?key=%s", c.baseURL, model, c.apiKey)
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("json marshal error: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("http request creation failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google API returned status %d: %s", resp.StatusCode, string(respBytes))
	}

	return respBytes, nil
}

func extractText(respBytes []byte) (string, error) {
	var resp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return "", fmt.Errorf("unmarshaling gemini response failed: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("gemini response contains no candidates")
	}

	return strings.TrimSpace(resp.Candidates[0].Content.Parts[0].Text), nil
}
