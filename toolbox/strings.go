package toolbox

import "sort"

func RemoveFromStings(texts []string, toRemove string) []string {
	stripped := []string{}
	for _, text := range texts {
		if text != toRemove {
			stripped = append(stripped, text)
		}
	}
	return stripped
}

func SortedStringsContains(list []string, search string) bool {
	idx := sort.SearchStrings(list, search)
	if idx >= len(list) {
		return false
	} else {
		return list[idx] == search
	}
}

func IndexOfString(slice []string, target string) int {
	for idx, s := range slice {
		if s == target {
			return idx
		}
	}
	return -1
}

func IsNilOrEmpty(text *string) bool {
	return text == nil || *text == ""
}
