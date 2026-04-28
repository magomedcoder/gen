package bitrix24server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"
)

type bitrixClient struct {
	baseURL      *url.URL
	httpClient   *http.Client
	logLevel     string
	retryMax     int
	retryBackoff time.Duration
}

var readOnlyBitrixMethods = map[string]struct{}{
	"tasks.task.list":          {},
	"tasks.task.get":           {},
	"task.commentitem.getlist": {},
}

func newBitrixClient(baseURL string, timeout time.Duration, logLevel string, retryMax int, retryBackoff time.Duration) (*bitrixClient, error) {
	trimmed := strings.TrimSpace(baseURL)
	if trimmed == "" {
		return nil, fmt.Errorf("bitrix webhook base is empty (set X-B24-Base in MCP headers)")
	}

	trimmed = strings.TrimSuffix(trimmed, "/")
	u, err := url.Parse(trimmed)
	if err != nil {
		return nil, fmt.Errorf("invalid bitrix webhook base: %w", err)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("bitrix webhook base must start with http:// or https://")
	}

	if u.Host == "" {
		return nil, fmt.Errorf("bitrix webhook base host is empty")
	}

	return &bitrixClient{
		baseURL: u,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		logLevel:     strings.ToLower(strings.TrimSpace(logLevel)),
		retryMax:     retryMax,
		retryBackoff: retryBackoff,
	}, nil
}

func (c *bitrixClient) call(ctx context.Context, method string, payload any) (map[string]any, error) {
	method = strings.TrimSpace(method)
	if method == "" {
		return nil, fmt.Errorf("method is empty")
	}
	if err := validateReadOnlyBitrixMethod(method); err != nil {
		return nil, err
	}
	reqID, _ := requestIDFromContext(ctx)

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	endpoint := *c.baseURL
	endpoint.Path = strings.TrimSuffix(endpoint.Path, "/") + "/" + method
	log.Printf("[b24-mcp] req_id=%s http method=%q url=%s payload_bytes=%d", reqID, method, endpoint.String(), len(body))
	if c.isDebugLogEnabled() {
		log.Printf("[b24-mcp] req_id=%s http method=%q request_body=%s", reqID, method, truncateForLog(prettyJSONString(maskJSONForLog(body)), 3000))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	var resp *http.Response
	for attempt := 0; ; attempt++ {
		resp, err = c.httpClient.Do(req)
		if err == nil {
			break
		}

		log.Printf("[b24-mcp] req_id=%s http method=%q request_err=%v attempt=%d", reqID, method, err, attempt+1)
		incHTTPError()
		if attempt >= c.retryMax {
			return nil, fmt.Errorf("request не удалось: %w", err)
		}
		time.Sleep(c.retryBackoff * time.Duration(attempt+1))
		req, _ = http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 5<<20))
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	raw = normalizeResponseEncoding(raw)
	incHTTPCall()
	log.Printf("[b24-mcp] req_id=%s http method=%q status=%d response_bytes=%d", reqID, method, resp.StatusCode, len(raw))
	if c.isDebugLogEnabled() {
		log.Printf("[b24-mcp] req_id=%s http method=%q response_body=%s", reqID, method, truncateForLog(prettyJSONString(maskJSONForLog(raw)), 3000))
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		log.Printf("[b24-mcp] req_id=%s http method=%q non_2xx status=%d body=%s", reqID, method, resp.StatusCode, truncateForLog(string(raw), 600))
		incHTTPError()
		if shouldRetryStatus(resp.StatusCode) && c.retryMax > 0 {
			for retry := 0; retry < c.retryMax; retry++ {
				time.Sleep(c.retryBackoff * time.Duration(retry+1))
				reqRetry, _ := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), bytes.NewReader(body))
				reqRetry.Header.Set("Content-Type", "application/json")
				respRetry, errRetry := c.httpClient.Do(reqRetry)
				if errRetry != nil {
					log.Printf("[b24-mcp] req_id=%s http method=%q retry_err=%v attempt=%d", reqID, method, errRetry, retry+1)
					incHTTPError()
					continue
				}

				retryRaw, _ := io.ReadAll(io.LimitReader(respRetry.Body, 5<<20))
				_ = respRetry.Body.Close()
				retryRaw = normalizeResponseEncoding(retryRaw)
				incHTTPCall()
				if respRetry.StatusCode >= http.StatusOK && respRetry.StatusCode < http.StatusMultipleChoices {
					var response map[string]any
					if err := json.Unmarshal(retryRaw, &response); err == nil {
						return response, nil
					}
				}
			}
		}
		return nil, wrapBitrixError(method, resp.StatusCode, raw, nil)
	}

	var response map[string]any
	if err := json.Unmarshal(raw, &response); err != nil {
		log.Printf("[b24-mcp] req_id=%s http method=%q decode_json_err=%v body=%s", reqID, method, err, truncateForLog(string(raw), 600))
		return nil, fmt.Errorf("decode json: %w", err)
	}

	if errObj, ok := response["error"]; ok {
		desc, _ := response["error_description"].(string)
		log.Printf("[b24-mcp] req_id=%s http method=%q api_error=%v desc=%q", reqID, method, errObj, desc)
		return nil, wrapBitrixError(method, resp.StatusCode, raw, response)
	}

	log.Printf("[b24-mcp] req_id=%s http method=%q ok", reqID, method)

	return response, nil
}

func validateReadOnlyBitrixMethod(method string) error {
	normalized := strings.ToLower(strings.TrimSpace(method))
	if normalized == "" {
		return fmt.Errorf("bitrix method is empty")
	}

	if _, ok := readOnlyBitrixMethods[normalized]; ok {
		return nil
	}

	return fmt.Errorf("bitrix method %q is blocked by read-only policy", method)
}

func (c *bitrixClient) isDebugLogEnabled() bool {
	return c.logLevel == "debug"
}

func shouldRetryStatus(status int) bool {
	return status == http.StatusTooManyRequests || status == http.StatusServiceUnavailable || status >= 500
}

func (c *bitrixClient) callTaskCommentItemGetList(ctx context.Context, taskID int, order map[string]any, filter map[string]any) (map[string]any, error) {
	payload := struct {
		TaskID int            `json:"TASKID"`
		Order  map[string]any `json:"ORDER,omitempty"`
		Filter map[string]any `json:"FILTER,omitempty"`
	}{
		TaskID: taskID,
		Order:  order,
		Filter: filter,
	}
	return c.call(ctx, "task.commentitem.getlist", payload)
}

func normalizeResponseEncoding(raw []byte) []byte {
	if len(raw) == 0 || utf8.Valid(raw) {
		return raw
	}

	decoded, err := charmap.Windows1251.NewDecoder().Bytes(raw)
	if err != nil || !utf8.Valid(decoded) || !json.Valid(decoded) {
		return raw
	}

	return decoded
}

func truncateForLog(s string, maxRunes int) string {
	if maxRunes <= 0 || utf8.RuneCountInString(s) <= maxRunes {
		return s
	}

	r := []rune(s)
	return string(r[:maxRunes]) + "...(truncated)"
}

func prettyJSONString(raw []byte) string {
	if len(raw) == 0 {
		return ""
	}

	var parsed any
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return string(raw)
	}

	formatted, err := json.MarshalIndent(parsed, "", "  ")
	if err != nil {
		return string(raw)
	}

	return string(formatted)
}

func maskJSONForLog(raw []byte) []byte {
	if len(raw) == 0 {
		return raw
	}

	var parsed any
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return raw
	}
	masked := maskValue(parsed)
	b, err := json.Marshal(masked)
	if err != nil {
		return raw
	}

	return b
}

func maskValue(v any) any {
	switch typed := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(typed))
		for k, vv := range typed {
			lk := strings.ToLower(strings.TrimSpace(k))
			if lk == "auth" || strings.Contains(lk, "token") || strings.Contains(lk, "password") || strings.Contains(lk, "secret") {
				out[k] = "***MASKED***"
				continue
			}
			out[k] = maskValue(vv)
		}
		return out
	case []any:
		out := make([]any, 0, len(typed))
		for _, item := range typed {
			out = append(out, maskValue(item))
		}
		return out
	default:
		return v
	}
}

func wrapBitrixError(method string, statusCode int, raw []byte, decoded map[string]any) error {
	var code string
	var desc string
	if decoded != nil {
		code, _ = decoded["error"].(string)
		desc, _ = decoded["error_description"].(string)
	} else {
		var m map[string]any
		if err := json.Unmarshal(raw, &m); err == nil {
			code, _ = m["error"].(string)
			desc, _ = m["error_description"].(string)
		}
	}
	hint := bitrixErrorHint(method, code, desc, statusCode)
	if code == "" && desc == "" {
		return fmt.Errorf("bitrix method %s не удалось with status %d. %s", method, statusCode, hint)
	}

	return fmt.Errorf("bitrix method %s не удалось: %s (%s). %s", method, code, strings.TrimSpace(desc), hint)
}

func bitrixErrorHint(method, code, desc string, statusCode int) string {
	msg := strings.ToLower(code + " " + desc)
	switch {
	case strings.Contains(msg, "wrong_arguments"), strings.Contains(msg, "expected to be of type"):
		return "Проверьте типы и названия параметров requestа."
	case strings.Contains(msg, "action_не удалось_to_be_processed"), strings.Contains(msg, "tasks_error_exception_#8"):
		if method == "task.commentitem.getlist" {
			return "Комментарий API может быть ограничен правами/новой карточкой задач. Попробуйте отключить include_comments."
		}
		return "Операция отклонена на стороне Bitrix24. Проверьте права и корректность полей."
	case strings.Contains(msg, "access denied"), statusCode == http.StatusForbidden:
		return "Недостаточно прав у пользователя/вебхука для этого метода."
	case statusCode == http.StatusTooManyRequests || strings.Contains(msg, "query_limit_exceeded"):
		return "Превышен лимит requestов Bitrix24, повторите позже."
	default:
		return "Проверьте параметры метода, права доступа и ограничения портала."
	}
}
