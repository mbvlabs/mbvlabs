package admin

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

func adminDateString(value sql.NullTime) string {
	if !value.Valid {
		return ""
	}

	return value.Time.Format("2006-01-02")
}

func commaSeparatedJSON(value string) json.RawMessage {
	items := []string{}
	for _, item := range strings.Split(value, ",") {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			items = append(items, trimmed)
		}
	}

	encoded, err := json.Marshal(items)
	if err != nil {
		return json.RawMessage("[]")
	}

	return encoded
}

func jsonArrayCSV(value json.RawMessage) string {
	var items []string
	if err := json.Unmarshal(value, &items); err != nil {
		return string(value)
	}

	return strings.Join(items, ", ")
}

func objectJSON(value string) (json.RawMessage, error) {
	if strings.TrimSpace(value) == "" {
		return json.RawMessage("{}"), nil
	}

	var obj map[string]any
	if err := json.Unmarshal([]byte(value), &obj); err != nil {
		return nil, fmt.Errorf("metadata must be valid JSON: %w", err)
	}

	encoded, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("metadata could not be encoded: %w", err)
	}

	return encoded, nil
}

func objectJSONString(value json.RawMessage) string {
	if len(value) == 0 {
		return "{}"
	}

	var obj map[string]any
	if err := json.Unmarshal(value, &obj); err != nil {
		return string(value)
	}

	encoded, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return string(value)
	}

	return string(encoded)
}
