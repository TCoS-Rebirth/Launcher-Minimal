package main

type Latest struct {
	Version  string `json:"version"`
	File     string `json:"file"`
	Checksum string `json:"checksum"`
	Server   string `json:"server"`
}

type Updates []struct {
	Update struct {
		AppliesTo  string `json:"applies_to"`
		Version    string `json:"version"`
		File       string `json:"file"`
		Patchnotes string `json:"patchnotes"`
		Checksum   string `json:"checksum"`
		Server     string `json:"server"`
		Enabled    bool   `json:"enabled"`
	} `json:"update"`
}

type Game struct {
	Path    string `json:"path"`
	Version string `json:"version"`
}
