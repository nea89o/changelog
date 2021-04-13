package pkg

import (
	"encoding/json"
	"os"
	"strings"
)

type Changes struct {
	Added   []string `json:"added"`
	Removed []string `json:"removed"`
	Changed []string `json:"changed"`
	Fixed   []string `json:"fixed"`
	Other   []string `json:"other"`
}
type Version struct {
	Changes
	Yanked bool   `json:"yanked"`
	Name   string `json:"name"`
}

type Changelog struct {
	Unreleased Changes   `json:"unreleased"`
	Past       []Version `json:"past"`
}

func NewEmptyChanges() Changes {
	return Changes{
		Added:   make([]string, 0),
		Removed: make([]string, 0),
		Changed: make([]string, 0),
		Fixed:   make([]string, 0),
		Other:   make([]string, 0),
	}
}

func NewVersion(name string) Version {
	return Version{
		Changes: NewEmptyChanges(),
		Name:    name,
		Yanked:  false,
	}
}

func NewChangelog() Changelog {
	return Changelog{
		Unreleased: NewEmptyChanges(),
		Past:       make([]Version, 0),
	}
}

func (v Version) Heading() (h string) {
	h = v.Name
	if v.Yanked {
		h += " [YANKED]"
	}
	return
}

func LoadOrCreateChangelog(path string) (Changelog, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return NewChangelog(), nil
	}
	return LoadChangelog(path)
}

func LoadChangelog(path string) (Changelog, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Changelog{}, err
	}
	var changelog Changelog
	err = json.Unmarshal(data, &changelog)
	if err != nil {
		return Changelog{}, err
	}
	return changelog, nil
}

func SaveChangelog(c Changelog, path string) (err error) {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return
	}
	err = os.WriteFile(path, data, 0644)
	return
}

func writeChangeSection(b *strings.Builder, changes []string, heading, verb string) {
	if len(changes) != 0 {
		b.WriteString("### ")
		b.WriteString(heading)
		b.WriteString("\n\n")
		for _, i := range changes {
			b.WriteString(" - ")
			b.WriteString(verb)
			b.WriteString(" ")
			b.WriteString(i)
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}
}

func (c Changes) writeChanges(b *strings.Builder, heading string) {
	if len(c.Added) != 0 || len(c.Removed) != 0 || len(c.Changed) != 0 || len(c.Fixed) != 0 || len(c.Other) != 0 {
		b.WriteString("## ")
		b.WriteString(heading)
		b.WriteString("\n\n")
		writeChangeSection(b, c.Added, "Additions", "Added")
		writeChangeSection(b, c.Removed, "Removals", "Removed")
		writeChangeSection(b, c.Changed, "Changes", "Changed")
		writeChangeSection(b, c.Fixed, "Fixes", "Fixed")
		writeChangeSection(b, c.Other, "Other", "")
		b.WriteString("\n")
	}
}

func (c Changelog) GenerateMarkdown() string {
	builder := strings.Builder{}
	c.Unreleased.writeChanges(&builder, "Unreleased")
	for _, version := range c.Past {
		version.writeChanges(&builder, version.Heading())
	}
	return builder.String()
}
