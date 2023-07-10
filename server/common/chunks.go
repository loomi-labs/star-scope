package common

func Chunks(text string, limit int) []string {
	if len(text) <= limit {
		return []string{text}
	}
	var result []string
	for len(text) > 0 {
		if len(text) < limit {
			limit = len(text)
		}
		result = append(result, text[:limit])
		text = text[limit:]
	}
	return result
}
