package str

import "fmt"

const reset string = "\033[0m"

// Red wraps the text in red color
func Red(text string) string {
	return fmt.Sprintf("\033[31m%s%s", text, reset)
}

// Green wraps the text in green color
func Green(text string) string {
	return fmt.Sprintf("\033[32m%s%s", text, reset)
}

// Yellow wraps the text in yellow color
func Yellow(text string) string {
	return fmt.Sprintf("\033[33m%s%s", text, reset)
}

// Blue wraps the text in blue color
func Blue(text string) string {
	return fmt.Sprintf("\033[34m%s%s", text, reset)
}

// Violet wraps the text in violet color
func Violet(text string) string {
	return fmt.Sprintf("\033[35m%s%s", text, reset)
}

// Cyan wraps the text in cyan color
func Cyan(text string) string {
	return fmt.Sprintf("\033[36m%s%s", text, reset)
}
