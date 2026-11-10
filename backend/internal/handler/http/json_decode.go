package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

const maxDecodedStringLength = 4096

func decodeJSONStrict(r *http.Request, dst interface{}) error {
	if r == nil || r.Body == nil {
		return errors.New("request body is required")
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		if errors.Is(err, io.EOF) {
			return errors.New("request body is required")
		}
		return fmt.Errorf("invalid request body: %w", err)
	}

	// Ensure single JSON object only.
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("request body must contain a single JSON object")
	}
	if err := validateDecodedPayloadStringLengths(dst, maxDecodedStringLength); err != nil {
		return err
	}
	return nil
}

func decodeJSONStrictAllowEmpty(r *http.Request, dst interface{}) error {
	if r == nil || r.Body == nil {
		return nil
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return fmt.Errorf("invalid request body: %w", err)
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("request body must contain a single JSON object")
	}
	if err := validateDecodedPayloadStringLengths(dst, maxDecodedStringLength); err != nil {
		return err
	}
	return nil
}

func validateDecodedPayloadStringLengths(payload interface{}, maxLen int) error {
	if payload == nil || maxLen <= 0 {
		return nil
	}
	return validateValueStringLengths(reflect.ValueOf(payload), maxLen, "$")
}

func validateValueStringLengths(v reflect.Value, maxLen int, path string) error {
	if !v.IsValid() {
		return nil
	}

	for v.Kind() == reflect.Interface || v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.String:
		if len([]rune(v.String())) > maxLen {
			return fmt.Errorf("input too long at %s (max %d chars)", path, maxLen)
		}
	case reflect.Struct:
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			if field.PkgPath != "" {
				continue
			}
			if err := validateValueStringLengths(v.Field(i), maxLen, path+"."+field.Name); err != nil {
				return err
			}
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if err := validateValueStringLengths(v.Index(i), maxLen, fmt.Sprintf("%s[%d]", path, i)); err != nil {
				return err
			}
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			k := key.Interface()
			childPath := fmt.Sprintf("%s[%v]", path, k)
			if key.Kind() == reflect.String && len([]rune(key.String())) > maxLen {
				return fmt.Errorf("input too long at %s key (max %d chars)", path, maxLen)
			}
			if err := validateValueStringLengths(v.MapIndex(key), maxLen, childPath); err != nil {
				return err
			}
		}
	}
	return nil
}
