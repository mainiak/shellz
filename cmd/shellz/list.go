package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/evilsocket/islazy/log"
	"github.com/evilsocket/shellz/plugins"

	"github.com/evilsocket/islazy/tui"
)

func showIdentsList() {
	rows := [][]string{}
	cols := []string{
		"Name",
		"Username",
		"Key",
		"Password",
		// "Path",
		"Comment",
	}

	keys := []string{}
	for k := range Idents {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, name := range keys {
		i := Idents[name]
		key := i.KeyFile
		if key == "" {
			key = tui.Dim("<empty>")
		}
		pass := strings.Repeat("*", len(i.Password))
		if pass == "" {
			pass = tui.Dim("<empty>")
		}
		comment := i.Comment
		if comment == "" {
			comment = tui.Dim("<empty>")
		}

		rows = append(rows, []string{
			tui.Bold(i.Name),
			i.Username,
			key,
			pass,
			// tui.Dim(i.Path),
			comment,
		})
	}

	fmt.Printf("\n%s\n", tui.Bold("identities"))
	tui.Table(os.Stdout, cols, rows)
}

func showPluginsList() {
	if plugins.Number() > 0 {
		rows := [][]string{}
		cols := []string{
			"Name",
			"Path",
		}

		plugins.Each(func(p *plugins.Plugin) {
			rows = append(rows, []string{
				tui.Bold(p.Name),
				tui.Dim(p.Path),
			})
		})

		fmt.Printf("\n%s\n", tui.Bold("plugins"))
		tui.Table(os.Stdout, cols, rows)
	}
}

func showShellsList() {
	if err, onShells = doShellSelection(onFilter, true); err != nil {
		log.Fatal("%s", err)
	} else if nShells = len(onShells); nShells == 0 {
		log.Fatal("no shell selected by the filter %s", tui.Dim(onFilter))
	}

	rows := [][]string{}
	cols := []string{
		"Name",
		"Groups",
		"Type",
		"Host",
		"Port",
		"Identity",
		"Enabled",
		"Comment",
	}

	hasTunnel := false
	for _, sh := range onShells {
		if !sh.Tunnel.Empty() {
			hasTunnel = true
			break
		}
	}

	if hasTunnel {
		cols = []string{
			"Name",
			"Groups",
			"Type",
			"Host",
			"Tunnel",
			"Port",
			"Identity",
			"Enabled",
			"Comment",
		}
	}

	keys := []string{}
	for k := range onShells {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, name := range keys {
		var row []string

		sh := onShells[name]
		en := tui.Green("✔")
		if !sh.Enabled {
			en = tui.Red("✖")
		}

		if hasTunnel {
			row = []string{
				tui.Bold(sh.Name),
				tui.Blue(strings.Join(sh.Groups, ", ")),
				tui.Dim(sh.Type),
				sh.Host,
				sh.Tunnel.String(),
				fmt.Sprintf("%d", sh.Port),
				tui.Yellow(sh.IdentityName),
				en,
				sh.Comment,
			}
		} else {
			row = []string{
				tui.Bold(sh.Name),
				tui.Blue(strings.Join(sh.Groups, ", ")),
				tui.Dim(sh.Type),
				sh.Host,
				fmt.Sprintf("%d", sh.Port),
				tui.Yellow(sh.IdentityName),
				en,
				sh.Comment,
			}
		}

		if !sh.Enabled {
			for i := range row {
				row[i] = tui.Dim(row[i])
			}
		}

		rows = append(rows, row)
	}

	fmt.Printf("\n%s\n", tui.Bold("shells"))
	tui.Table(os.Stdout, cols, rows)
}

func showList() {
	showIdentsList()
	showPluginsList()
	showShellsList()
}
