package workspaceappearance

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"regexp"
	"strconv"
)

var hexColor = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`) // Opaque RGB colors only.

var defaultColors = map[string]map[string]string{
	"light": {"background": "#F8F2E8", "surface": "#EFE4D5", "foreground": "#2B201C", "muted": "#5A4941", "border": "#D8C8B6", "primary": "#A13224"},
	"dark":  {"background": "#211B18", "surface": "#2D2520", "foreground": "#F2E9DC", "muted": "#D0C1B1", "border": "#514237", "primary": "#E17861"},
}

func parseObject(value any) (map[string]any, error) {
	object, ok := value.(map[string]any)
	if !ok || object == nil {
		return nil, errors.New("expected object")
	}
	return object, nil
}

// ValidateOverrides rejects unsupported values before persistence or CSS use.
func ValidateOverrides(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	var raw any
	if err := decoder.Decode(&raw); err != nil {
		return fmt.Errorf("decode appearance: %w", err)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return errors.New("unexpected trailing JSON")
	}
	root, err := parseObject(raw)
	if err != nil {
		return err
	}
	for key, value := range root {
		switch key {
		case "colors":
			err = validateColors(value)
		case "typography":
			err = validateTypography(value)
		default:
			return fmt.Errorf("unknown appearance field %q", key)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func validateColors(raw any) error {
	modes, err := parseObject(raw)
	if err != nil {
		return err
	}
	for mode, value := range modes {
		defaults, ok := defaultColors[mode]
		if !ok {
			return errors.New("unknown color mode")
		}
		fields, err := parseObject(value)
		if err != nil {
			return err
		}
		resolved := make(map[string]string, len(defaults))
		for key, color := range defaults {
			resolved[key] = color
		}
		for key, color := range fields {
			text, ok := color.(string)
			if _, known := defaults[key]; !known || !ok || !hexColor.MatchString(text) {
				return errors.New("invalid color token")
			}
			resolved[key] = text
		}
		for _, surface := range []string{"background", "surface"} {
			if contrast(resolved["foreground"], resolved[surface]) < 4.5 || contrast(resolved["muted"], resolved[surface]) < 4.5 {
				return errors.New("insufficient text contrast")
			}
		}
		onPrimary := "#FFFFFF"
		if mode == "dark" {
			onPrimary = "#211B18"
		}
		if contrast(resolved["primary"], onPrimary) < 4.5 {
			return errors.New("insufficient primary contrast")
		}
	}
	return nil
}

func contrast(a, b string) float64 {
	x, y := luminance(a), luminance(b)
	if x > y {
		x, y = y, x
	}
	return (y + .05) / (x + .05)
}

func luminance(color string) float64 {
	weights := []float64{.2126, .7152, .0722}
	result := 0.0
	for i := range 3 {
		value, _ := strconv.ParseUint(color[1+i*2:3+i*2], 16, 8)
		channel := float64(value) / 255
		if channel <= .04045 {
			channel /= 12.92
		} else {
			channel = math.Pow((channel+.055)/1.055, 2.4)
		}
		result += channel * weights[i]
	}
	return result
}

func validateTypography(raw any) error {
	fields, err := parseObject(raw)
	if err != nil {
		return err
	}
	for key, value := range fields {
		switch key {
		case "english":
			if value != "manrope" && value != "system" {
				return errors.New("unsupported English font")
			}
		case "arabic":
			if value != "ibm-plex-arabic" && value != "system" {
				return errors.New("unsupported Arabic font")
			}
		case "heading":
			if value != "fraunces" && value != "manrope" {
				return errors.New("unsupported heading font")
			}
		case "baseSize":
			n, ok := value.(float64)
			if !ok || n < 14 || n > 18 || math.Trunc(n) != n {
				return errors.New("invalid base size")
			}
		default:
			return errors.New("unknown typography field")
		}
	}
	return nil
}
