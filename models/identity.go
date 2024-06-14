package models

import (
	"encoding/json"
	"io"
	"os"
)

type Identity struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	KeyFile  string `json:"key"`
	Password string `json:"password"`
	Path     string `json:"-"`
	Comment  string `json:"-"`
}

func LoadIdent(path string) (err error, ident Identity) {
	ident = Identity{Path: path}
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	raw, err := io.ReadAll(file)
	if err != nil {
		return
	}

	if err = json.Unmarshal(raw, &ident); err != nil {
		return
	}

	return
}
