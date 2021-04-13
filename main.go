package main

import (
	"fmt"
	"os"

	changelog "github.com/romangraef/changelog/pkg"
)

func main() {
	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: CHANGELOG <add/remove/change/fix/other> <what>")
		os.Exit(1)
	}
	path := "changelog.json"
	action := args[0]
	what := args[1]
	c, err := changelog.LoadOrCreateChangelog(path)
	errorIf(err, "Failed to load changelog.json")
	switch action {
	case "add":
		c.Unreleased.Added = append(c.Unreleased.Added, what)
	case "remove":
		c.Unreleased.Removed = append(c.Unreleased.Removed, what)
	case "change":
		c.Unreleased.Changed = append(c.Unreleased.Changed, what)
	case "fix":
		c.Unreleased.Fixed = append(c.Unreleased.Fixed, what)
	case "other":
		c.Unreleased.Other = append(c.Unreleased.Other, what)
	case "release":
		c.Past = append(c.Past, changelog.Version{
			Changes: c.Unreleased,
			Name:    what,
			Yanked:  false,
		})
		c.Unreleased = changelog.NewEmptyChanges()
	case "write":
		os.WriteFile(what, []byte(c.GenerateMarkdown()), 0644)
	default:
		fmt.Fprintf(os.Stderr, "ERROR: Unknown action %v\n", action)
	}
	errorIf(changelog.SaveChangelog(c, path), "Failed to save changelog.json")
}

func errorIf(err error, msg string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v: %v\n", msg, err)
		os.Exit(1)
	}
}
