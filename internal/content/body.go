package content

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Body struct {
	images   []*Image
	sections []*Section
}

func New() *Body {
	return &Body{}
}

func (b *Body) AddSection(x, y, maxLength int) *Section {
	section := &Section{
		x:         x,
		y:         y,
		maxLength: maxLength,
	}
	b.sections = append(b.sections, section)
	return section
}

func (b *Body) AddAsciiArt(x, y int, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" class="ascii">`+"\n", x, y))
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		sb.WriteString(tspanPos(x, y, "", line))
		sb.WriteString("\n")
		y += 20
	}
	sb.WriteString("</text>\n")

	b.images = append(b.images, &Image{
		image: sb.String(),
	})

	return nil
}

func (b *Body) String() string {
	var sb strings.Builder

	for _, i := range b.images {
		sb.WriteString(i.String())
	}

	for _, s := range b.sections {
		sb.WriteString(s.String())
	}

	return sb.String()
}

func (b *Body) WriteTemplate(templates ...string) error {
	var files [][]byte

	for i, t := range templates {
		file, err := os.ReadFile(t)
		if err != nil {
			return err
		}
		files = append(files, file)

		// fix names for writing
		tmpArr := strings.Split(t, "/")
		tmpStr := tmpArr[len(tmpArr)-1]
		tmpArr = strings.Split(tmpStr, ".")
		templates[i] = tmpArr[len(tmpArr)-2] + ".svg"
	}

	for i, f := range files {
		str := string(f)
		str = strings.ReplaceAll(str, "{{content}}", b.String())

		err := os.WriteFile(templates[i], []byte(str), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}
