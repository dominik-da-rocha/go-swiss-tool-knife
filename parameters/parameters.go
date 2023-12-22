package parameters

import "go-tool-box/strings"

func ToInt64(value *int64, def int64) int64 {
	if value == nil {
		return def
	} else {
		return *value
	}
}

func ToString(value *string, def string) string {
	if value == nil {
		return def
	} else {
		return *value
	}
}

func ToAnyStringOf(text *string, allowed []string, def string) string {
	if strings.IsNilOrEmpty(text) {
		return def
	}
	idx := strings.IndexOfString(allowed, *text)
	if idx < 0 {
		return def
	}
	return *text
}
