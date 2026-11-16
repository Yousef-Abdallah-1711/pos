package csvutil

import (
	"strings"
	"unicode"
)

// SafeTextCell prefixes spreadsheet formula-leading text so spreadsheet apps
// treat it as a value. Callers must use it only for text cells; numeric cells
// should retain their native CSV representation.
func SafeTextCell(value string) string {
	for _, character := range value {
		if unicode.IsSpace(character) || unicode.IsControl(character) {
			continue
		}
		if strings.ContainsRune("=+-@", character) {
			return "'" + value
		}
		return value
	}
	return value
}
