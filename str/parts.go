package str

import "strings"

// CleanSplit splits the text by separator, and trims each part's extra whitespace
func CleanSplit(text, sep string) []string {
	return CleanSplitter(sep)(text)
}

// CleanSplitter creates a function that splits the text by separator, and trims each part's extra whitespace
func CleanSplitter(sep string) func(string) []string {
	return func(text string) []string {
		parts := strings.Split(text, sep)
		for i, part := range parts {
			parts[i] = strings.TrimSpace(part)
		}
		return parts
	}
}

// CleanSplitN splits the text by separator, maximum of N parts, then trims each part's extra whitespace
func CleanSplitN(text, sep string, n int) []string {
	return CleanSplitterN(sep, n)(text)
}

// CleanSplitterN creates a function that splits the text by separator, maximum of N parts, and trims each part's whitespace
func CleanSplitterN(sep string, n int) func(string) []string {
	return func(text string) []string {
		parts := strings.SplitN(text, sep, n)
		for i, part := range parts {
			parts[i] = strings.TrimSpace(part)
		}
		return parts
	}
}

// SpaceSplit splits the text by whitespace
func SpaceSplit(text string) []string {
	return strings.Fields(strings.TrimSpace(text))
}

// Join joins the strings parts by given glue
func Join(glue string, parts ...string) string {
	return strings.Join(parts, glue)
}
