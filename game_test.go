package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// TestIsInstalled tests the isInstalled function
func TestIsInstalled(t *testing.T) {
	// Create temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "game_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Save original gameInfo path and restore after test
	originalGameInfo := gameInfo
	defer func() { gameInfo = originalGameInfo }()
	gameInfo = filepath.Join(tmpDir, "gameInfo.json")

	tests := []struct {
		name     string
		gameData *Game
		want     bool
	}{
		{
			name: "valid installation",
			gameData: &Game{
				Path:    "test/path",
				Version: "1.0.0",
			},
			want: true,
		},
		{
			name: "missing version",
			gameData: &Game{
				Path: "test/path",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.gameData != nil {
				data, _ := json.Marshal(tt.gameData)
				err := os.WriteFile(gameInfo, data, 0644)
				if err != nil {
					t.Fatal(err)
				}
			}

			if got := isInstalled(); got != tt.want {
				t.Errorf("isInstalled() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestFetchLatestVersion tests the fetchLatestVersion function
func TestFetchLatestVersion(t *testing.T) {
	// Create test server
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/latest.json" {
			t.Errorf("Expected path /latest.json, got %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}

		latest := Latest{
			Version:  "1.0.0",
			File:     "game.zip",
			Checksum: "abc123",
			Server:   "test-server",
		}
		json.NewEncoder(w).Encode(latest)
	}))
	defer testServer.Close()

	// Save original fileServer and restore after test
	originalFileServer := fileServer
	defer func() { fileServer = originalFileServer }()
	fileServer = testServer.URL + "/"

	latest, err := fetchLatestVersion()
	if err != nil {
		t.Fatalf("fetchLatestVersion() error = %v", err)
	}

	if latest.Version != "1.0.0" {
		t.Errorf("Expected version 1.0.0, got %s", latest.Version)
	}
}

// TestUpdateVersion tests the updateVersion function
func TestUpdateVersion(t *testing.T) {
	// Create temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "game_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Save original gameInfo path and restore after test
	originalGameInfo := gameInfo
	defer func() { gameInfo = originalGameInfo }()
	gameInfo = filepath.Join(tmpDir, "gameInfo.json")

	// Create initial game info
	initialGame := Game{
		Path:    "test/path",
		Version: "1.0.0",
	}
	data, _ := json.Marshal(initialGame)
	err = os.WriteFile(gameInfo, data, 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Test updating version
	err = updateVersion("2.0.0")
	if err != nil {
		t.Fatalf("updateVersion() error = %v", err)
	}

	// Verify the update
	file, err := os.Open(gameInfo)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	var updatedGame Game
	bytes, _ := io.ReadAll(file)
	err = json.Unmarshal(bytes, &updatedGame)
	if err != nil {
		t.Fatal(err)
	}

	if updatedGame.Version != "2.0.0" {
		t.Errorf("Expected version 2.0.0, got %s", updatedGame.Version)
	}
}

// TestGetGameInfo tests the getGameInfo function
func TestGetGameInfo(t *testing.T) {
	// Create temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "game_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Save original gameInfo path and restore after test
	originalGameInfo := gameInfo
	defer func() { gameInfo = originalGameInfo }()
	gameInfo = filepath.Join(tmpDir, "gameInfo.json")

	// Test cases
	tests := []struct {
		name     string
		gameData *Game
		wantErr  bool
	}{
		{
			name: "valid game info",
			gameData: &Game{
				Path:          "test/path",
				Version:       "1.0.0",
				KeepDownloads: true,
				NoLaunch:      false,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.gameData != nil {
				data, _ := json.Marshal(tt.gameData)
				err := os.WriteFile(gameInfo, data, 0644)
				if err != nil {
					t.Fatal(err)
				}
			}

			got, err := getGameInfo()
			if (err != nil) != tt.wantErr {
				t.Errorf("getGameInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.Version != tt.gameData.Version {
					t.Errorf("getGameInfo() Version = %v, want %v", got.Version, tt.gameData.Version)
				}
				if got.Path != tt.gameData.Path {
					t.Errorf("getGameInfo() Path = %v, want %v", got.Path, tt.gameData.Path)
				}
			}
		})
	}
}
