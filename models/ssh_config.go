package models

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/evilsocket/islazy/log"
	"github.com/kevinburke/ssh_config"
)

/*
 * https://linux.die.net/man/5/ssh_config
 */

const (
	SSH_CONFIG_HOSTNAME = "HostName"
	SSH_CONFIG_IDENTITY = "IdentityFile"
	SSH_CONFIG_PORT     = "Port"
	SSH_CONFIG_USER     = "User"
)

func parseSSHConfigEntry(nodes []ssh_config.Node) (string, string, string, string) {
	var hostname, identity, port, user string

	for _, node := range nodes {
		// Manipulate the nodes as you see fit, or use a type switch to
		// distinguish between Empty, KV, and Include nodes.
		log.Debug("> %s", node.String())
		s := strings.TrimSpace(node.String())

		// ignore comments
		if strings.HasPrefix(s, "#") {
			continue
		}

		log.Debug(">> %s", s)
		parts := strings.Split(s, " ")
		log.Debug(">> %s", parts)

		switch parts[0] {
		case SSH_CONFIG_HOSTNAME:
			hostname = parts[1]
		case SSH_CONFIG_IDENTITY:
			identity = parts[1]
		case SSH_CONFIG_PORT:
			port = parts[1]
		case SSH_CONFIG_USER:
			user = parts[1]
		default:
			continue // probably redundant
		}
	}

	return hostname, identity, port, user
}

func loadSSHConfig(idents Identities, shells Shells, groups Groups) error {
	f, _ := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "config"))
	cfg, _ := ssh_config.Decode(f)

	for idx, host := range cfg.Hosts {
		log.Debug("patterns:", host.Patterns)

		// ignore if multiple host patterns (Host line separated by ',')
		if len(host.Patterns) != 1 {
			continue
		}

		// Skip glob patterns
		if strings.Contains(host.Patterns[0].String(), "*") {
			continue
		}

		shell_name := fmt.Sprintf("@host%02d", idx)
		identity_name := fmt.Sprintf("@ident%02d", idx)

		shell := NewShell("")
		shell.Name = shell_name
		shell.Comment = host.Patterns[0].String()

		cfg_hostname, cfg_identity, cfg_port, cfg_user := parseSSHConfigEntry(host.Nodes)

		if cfg_hostname != "" {
			shell.Host = cfg_hostname
		} else {
			shell.Host = shell.Comment
		}

		if cfg_port != "" {
			port_number, err := strconv.Atoi(cfg_port)
			if err != nil {
				log.Error("Invalid port number: %s", cfg_port)
				return err
			}
			shell.Port = port_number
		}

		create_identity := false
		ident := Identity{
			Name: identity_name,
		}

		if cfg_identity != "" {
			create_identity = true
			ident.KeyFile = "@agent"
			ident.Comment = cfg_identity
		}

		if cfg_user != "" {
			create_identity = true
			ident.Username = cfg_user
		}

		if create_identity {
			idents[ident.Name] = ident
			shell.IdentityName = identity_name
			shell.Identity = &ident
		}

		shells[shell.Name] = shell
	}

	// Print the config to stdout:
	log.Debug(cfg.String())

	return nil
}
