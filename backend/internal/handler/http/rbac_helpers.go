package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"pos-saas/internal/domain"
)

type permissionInput struct {
	ModuleID int64
	Level    string
}

func decodeJSONMap(r *http.Request) (map[string]interface{}, error) {
	decoder := json.NewDecoder(r.Body)
	decoder.UseNumber()
	var payload map[string]interface{}
	if err := decoder.Decode(&payload); err != nil {
		return nil, err
	}
	return payload, nil
}

func getStringField(payload map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		value, ok := payload[key]
		if !ok || value == nil {
			continue
		}
		if s, ok := value.(string); ok && strings.TrimSpace(s) != "" {
			return strings.TrimSpace(s)
		}
	}
	return ""
}

func getInt64Field(payload map[string]interface{}, keys ...string) (int64, bool, error) {
	for _, key := range keys {
		value, ok := payload[key]
		if !ok || value == nil {
			continue
		}
		id, err := parseInterfaceToInt64(value)
		if err != nil {
			return 0, true, err
		}
		return id, true, nil
	}
	return 0, false, nil
}

func parseInterfaceToInt64(value interface{}) (int64, error) {
	switch v := value.(type) {
	case json.Number:
		return v.Int64()
	case float64:
		return int64(v), nil
	case float32:
		return int64(v), nil
	case int:
		return int64(v), nil
	case int64:
		return v, nil
	case string:
		if strings.TrimSpace(v) == "" {
			return 0, errors.New("empty string for numeric field")
		}
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid numeric string: %w", err)
		}
		return parsed, nil
	default:
		return 0, fmt.Errorf("unsupported numeric type: %T", value)
	}
}

func normalizePermissionLevel(level string) (string, error) {
	normalized := strings.ToUpper(strings.TrimSpace(level))
	if normalized == "" {
		return "", errors.New("permission level is required")
	}
	if _, ok := domain.PermissionValues[normalized]; !ok {
		return "", fmt.Errorf("invalid permission level: %s", level)
	}
	return normalized, nil
}

func parsePermissionInputs(raw interface{}) ([]permissionInput, error) {
	if raw == nil {
		return nil, nil
	}

	items, ok := raw.([]interface{})
	if !ok {
		return nil, errors.New("permissions must be an array")
	}

	permissions := make([]permissionInput, 0, len(items))
	for _, item := range items {
		entry, ok := item.(map[string]interface{})
		if !ok {
			return nil, errors.New("invalid permission entry")
		}

		moduleID, ok, err := getInt64Field(entry, "module_id", "moduleId")
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, errors.New("module_id is required")
		}

		levelRaw := getStringField(entry, "permission_level", "level")
		level, err := normalizePermissionLevel(levelRaw)
		if err != nil {
			return nil, err
		}

		permissions = append(permissions, permissionInput{
			ModuleID: moduleID,
			Level:    level,
		})
	}

	return permissions, nil
}
