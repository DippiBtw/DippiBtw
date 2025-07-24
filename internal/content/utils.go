package content

import (
	"fmt"
	"strings"
	"time"
)

const prefix = "-"
const suffix = "-—-"
const repeat = "—"

func addTitleLine(title string, maxLength int) string {
	line := strings.Repeat(repeat, maxLength-len(title)-3)
	return fmt.Sprintf(` <tspan class="title">%s</tspan>`, transform(line, prefix, suffix))
}

func transform(s, first, last3 string) string {
	runes := []rune(s) // Handle Unicode safely

	if len(runes) < 4 {
		// Not enough runes to replace first + last 3
		return s
	}

	return first + string(runes[1:len(runes)-3]) + last3
}

func tspanPos(x, y int, class, content string) string {
	if class != "" {
		class = fmt.Sprintf(` class="%s"`, class)
	}
	return fmt.Sprintf(`<tspan x="%d" y="%d"%s>%s</tspan>`, x, y, class, content)
}

func tspan(class, content string) string {
	if class != "" {
		class = fmt.Sprintf(` class="%s"`, class)
	}
	return fmt.Sprintf(`<tspan%s>%s</tspan>`, class, content)
}

func addDots(key string, valueLen, maxLength int) string {
	// Remove three for the extra characters ". " and ":" when creating the entry text
	repeat := maxLength - len(key) - valueLen - 5
	if repeat < 0 {
		return ""
	}
	return fmt.Sprintf(" %s ", strings.Repeat(".", repeat))
}

func FormatDuration(d time.Duration) string {
	var years, months, days int
	totalDays := int(d.Hours() / 24)

	years = totalDays / 365
	remainingDays := totalDays % 365

	months = remainingDays / 30
	days = remainingDays % 30

	return fmt.Sprintf("%d years, %d months, %d days", years, months, days)
}
