package models

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const (
	defaultEnabled   = true
	defaultHost      = "localhost"
	defaultPort      = 22
	defaultIdentity  = "default"
	defaultType      = "ssh"
	defaultHTTPS     = true
	defaultInsecure  = false
	defaultNamespace = "default"
	defaultContainer = ""
	defaultPod       = ""
	defaultComment   = ""
)

type Shell struct {
	Name         string    `json:"name"`
	Host         string    `json:"host"`
	Port         int       `json:"port"`
	IdentityName string    `json:"identity"`
	Type         string    `json:"type"`
	Ciphers      []string  `json:"ciphers"`
	HTTPS        bool      `json:"https"`
	Insecure     bool      `json:"insecure"`
	Enabled      bool      `json:"enabled"`
	Groups       []string  `json:"groups"`
	Proxy        Proxy     `json:"proxy"`
	Tunnel       Tunnel    `json:"tunnel"`
	Namespace    string    `json:"namespace"`
	Container    string    `json:"container"`
	Pod          string    `json:"pod"`
	Identity     *Identity `json:"-"`
	Path         string    `json:"-"`
	Comment      string    `json:"-"`
}

func LoadShell(path string, idents Identities) (err error, shell Shell) {
	shell = Shell{
		Enabled:      defaultEnabled,
		Path:         path,
		Host:         defaultHost,
		Port:         defaultPort,
		Type:         defaultType,
		IdentityName: defaultIdentity,
		HTTPS:        defaultHTTPS,
		Insecure:     defaultInsecure,
		Namespace:    defaultNamespace,
		Container:    defaultContainer,
		Pod:          defaultPod,
		Comment:      defaultComment,
	}

	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	raw, err := io.ReadAll(file)
	if err != nil {
		return
	}

	if err = json.Unmarshal(raw, &shell); err != nil {
		return fmt.Errorf("error decoding '%s': %s", path, err), shell
	} else if ident, found := idents[shell.IdentityName]; !found {
		return fmt.Errorf("shell '%s' referenced an unknown identity '%s'", path, shell.IdentityName), shell
	} else {
		shell.Identity = &ident
	}

	return
}

func (sh Shell) Save() error {
	if data, err := json.MarshalIndent(sh, "", "  "); err != nil {
		return err
	} else if err = os.WriteFile(sh.Path, data, 0644); err != nil {
		return err
	}
	return nil
}
