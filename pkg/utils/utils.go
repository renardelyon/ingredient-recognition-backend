package utils

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

func ExtractJSONFromString(s string) (string, error) {
	text := strings.TrimSpace(s)

	if IsValidJSON(text) {
		return text, nil
	}

	codeBlockRegex := regexp.MustCompile("(?s)```(?:json)?\\s*(.+?)```")
	matches := codeBlockRegex.FindStringSubmatch(s)
	if len(matches) > 1 {
		extracted := strings.TrimSpace(matches[1])
		if IsValidJSON(extracted) {
			return extracted, nil
		}
	}

	if strings.Contains(text, "{") {
		start := strings.Index(text, "{")
		end := strings.LastIndex(text, "}")
		if start != -1 && end != -1 && end > start {
			extracted := text[start : end+1]
			if IsValidJSON(extracted) {
				return extracted, nil
			}
		}
	}

	if strings.Contains(text, "[") {
		start := strings.Index(text, "[")
		end := strings.LastIndex(text, "]")
		if start != -1 && end != -1 && end > start {
			extracted := text[start : end+1]
			if IsValidJSON(extracted) {
				return extracted, nil
			}
		}
	}

	return "", fmt.Errorf("malformed JSON in response")
}

func IsValidJSON(s string) bool {
	var js any
	return json.Unmarshal([]byte(s), &js) == nil
}
