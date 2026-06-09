package recaptcha

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const verifyURL = "https://www.google.com/recaptcha/api/siteverify"

type Verifier struct {
	secret     string
	httpClient *http.Client
}

type Response struct {
	Success     bool     `json:"success"`
	Score       float64  `json:"score"`
	Action      string   `json:"action"`
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes"`
}

func NewVerifier(secret string) *Verifier {
	return &Verifier{
		secret: secret,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (v *Verifier) Verify(token, expectedAction string, minScore float64) (*Response, error) {
	if v.secret == "" {
		return &Response{Success: true}, nil
	}

	data := url.Values{}
	data.Set("secret", v.secret)
	data.Set("response", token)

	resp, err := v.httpClient.Post(verifyURL, "application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("recaptcha verify: %w", err)
	}
	defer resp.Body.Close()

	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("recaptcha decode: %w", err)
	}

	if !result.Success {
		return &result, fmt.Errorf("recaptcha verification failed: %v", result.ErrorCodes)
	}

	if expectedAction != "" && result.Action != expectedAction {
		return &result, fmt.Errorf("recaptcha action mismatch: expected %q, got %q", expectedAction, result.Action)
	}

	if result.Score < minScore {
		return &result, fmt.Errorf("recaptcha score too low: %.2f < %.2f", result.Score, minScore)
	}

	return &result, nil
}
