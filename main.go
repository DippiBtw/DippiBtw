package main

import (
	"context"
	"log"
	"runtime"
	"strings"
	"time"

	"github.com/DippiBtw/DippiBtw/internal/content"
	"github.com/DippiBtw/DippiBtw/internal/queries"
	"github.com/DippiBtw/DippiBtw/internal/yamlHandler"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Inspired by [@Andrew6rant]'s profile

func main() {
	cfg, err := yamlHandler.LoadYAML("assets/info.yaml")
	if err != nil {
		log.Fatal(err)
	}

	body := content.New()

	body.AddAsciiArt(50, 45, "assets/totoro.txt")
	body.AddAsciiArt(650, 335, "assets/cursed.txt")

	section1 := body.AddSection(390, 30, 60)
	me := section1.AddParagraph(cfg.Section1[0].Title)
	uptime, err := time.Parse(time.RFC3339, cfg.Section1[0].Custom.Uptime)
	me.AddEntry("age", content.FormatDuration(time.Since(uptime)))
	addMapOfArrays(me, "hobbies", cfg.Section1[0].Content.Hobbies)
	addMapOfArrays(me, "languages", cfg.Section1[0].Content.Languages)
	addMapOfArrays(me, "system", cfg.Section1[0].Content.System)

	section2 := body.AddSection(30, 340, 65)
	contacts := section2.AddParagraph(cfg.Section2[0].Title)
	addMapOfStrings(contacts, "email", cfg.Section2[0].Content.Email)
	addMapOfStrings(contacts, "socials", cfg.Section2[0].Content.Socials)

	ctx := context.Background()
	info, err := queries.GetInfo(ctx)
	if err != nil {
		log.Fatal(err)
	}

	p := message.NewPrinter(language.English)
	stats := section2.AddParagraph("github.stats")
	stats.AddEntry("stat.commits", p.Sprintf("%d", info.Commits))
	stats.AddEntry("stat.followers", p.Sprintf("%d", info.Followers))
	stats.AddEntry("stat.lines.of.code", p.Sprintf(`%d ( [%d++], {%d--} )`, info.Additions-info.Deletions, info.Additions, info.Deletions))
	stats.AddEntry("stat.repos", p.Sprintf("%d { Contributed: %d }", info.Repos, info.ContributedTo))
	stats.AddEntry("stat.stars", p.Sprintf("%d", info.Stars))

	section2.AddParagraph("Updated " + time.Now().Format("Mon, Jan 2, 2006 at 15:04:05 MST"))

	templates := []string{"assets/dark.xhtml", "assets/light.xhtml"}
	var inkscape string

	if runtime.GOOS == "windows" {
		inkscape = `C:\Program Files\Inkscape\bin\inkscape`
	}

	err = body.ConvertSvgToPng(inkscape, templates...)
	if err != nil {
		log.Fatal(err)
	}
}

func addMapOfArrays(part *content.Part, prefix string, m map[string][]string) {
	var str string
	for k, v := range m {
		for _, s := range v {
			str += s + ", "
		}

		part.AddEntry(prefix+"."+k, strings.TrimSuffix(str, ", "))
		str = ""
	}
}

func addMapOfStrings(part *content.Part, prefix string, m map[string]string) {
	for k, v := range m {
		part.AddEntry(prefix+"."+k, v)
	}
}
