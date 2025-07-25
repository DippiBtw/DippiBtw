package content

import (
	"fmt"
	"sort"
	"strings"
)

type Section struct {
	parts     []*Part
	x         int
	y         int
	maxLength int
}

func (s *Section) AddParagraph(title string) *Part {
	part := &Part{
		title:   title,
		entries: make(map[string]string),
	}
	s.parts = append(s.parts, part)
	return part
}

func (s *Section) String() string {
	var sb strings.Builder

	x := s.x
	y := s.y

	sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" fill="#c9d1d9">`+"\n", s.x, s.y))
	for _, p := range s.parts {
		sb.WriteString(s.addTitle(p.title))

		var keys []string
		for key, _ := range p.entries {
			keys = append(keys, key)
		}

		sort.Strings(keys)
		var previous string
		for _, key := range keys {
			if previous != "" && previous != strings.Split(key, ".")[0] {
				sb.WriteString(s.addEntry("", ""))
			}
			previous = strings.Split(key, ".")[0]
			sb.WriteString(s.addEntry(key, p.entries[key]))
		}

		s.y += 20
	}
	sb.WriteString(fmt.Sprintf("</text>\n"))

	s.x = x
	s.y = y

	return sb.String()
}

func (s *Section) addTitle(title string) string {
	ret := tspanPos(s.x, s.y, "title", fmt.Sprintf("- %s %s", title, addTitleLine(title, s.maxLength)))
	s.y += 20
	return ret
}

func (s *Section) addEntry(key, value string) string {
	length := len(value)
	if strings.Contains(value, "[") {
		length = len(value) - 4
		value = strings.ReplaceAll(value, "[", `<tspan class="addColor">`)
		value = strings.ReplaceAll(value, "]", `</tspan>`)
		value = strings.ReplaceAll(value, "{", `<tspan class="delColor">`)
		value = strings.ReplaceAll(value, "}", `</tspan>`)
	}

	ret := fmt.Sprintf("%s%s:%s%s\n",
		tspanPos(s.x, s.y, "cc", ". "),
		tspan("key", key),
		tspan("cc", addDots(key, length, s.maxLength)),
		tspan("value", value))
	s.y += 20

	return ret
}
