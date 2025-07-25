package content

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/DippiBtw/DippiBtw/internal/inkscape"
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
		y += 17
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
	const filePerm = 0o644

	for _, t := range templates {
		contentBytes, err := os.ReadFile(t)
		if err != nil {
			return fmt.Errorf("reading template %q: %w", t, err)
		}

		outputFile := filepath.Clean(getOutput(t, ".svg"))

		// Do not change input contentBytes; create a new string for replaced content
		contentStr := string(contentBytes)
		replaced := strings.ReplaceAll(contentStr, "{{content}}", b.String())

		if err := os.WriteFile(outputFile, []byte(replaced), filePerm); err != nil {
			return fmt.Errorf("writing output file %q: %w", outputFile, err)
		}
	}

	return nil
}

// REQUIRES INKSCAPE
func (b *Body) ConvertSvgToPng(inkscapePath string, templates ...string) error {
	const filePerm = 0o644

	conv := inkscape.New()
	if inkscapePath != "" {
		conv.SetBinary(inkscapePath)
	}

	for _, t := range templates {
		contentBytes, err := os.ReadFile(t)
		if err != nil {
			return fmt.Errorf("reading template %q: %w", t, err)
		}

		outputFile := filepath.Clean(getOutput(t, ".png"))

		contentStr := string(contentBytes)
		replaced := strings.ReplaceAll(contentStr, "{{content}}", b.String())

		pngData, err := conv.Convert([]byte(replaced))
		if err != nil {
			return fmt.Errorf("converting SVG to PNG for %q: %w", t, err)
		}

		if err := os.WriteFile(outputFile, pngData, filePerm); err != nil {
			return fmt.Errorf("writing PNG file %q: %w", outputFile, err)
		}
	}

	return nil
}

func getOutput(str, suffix string) string {
	tmpArr := strings.Split(str, "/")
	tmpStr := tmpArr[len(tmpArr)-1]
	tmpArr = strings.Split(tmpStr, ".")
	return "output/" + tmpArr[len(tmpArr)-2] + suffix
}
