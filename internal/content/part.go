package content

type Part struct {
	title   string
	entries map[string]string
}

func (s *Part) AddEntry(key, value string) {
	s.entries[key] = value
}
