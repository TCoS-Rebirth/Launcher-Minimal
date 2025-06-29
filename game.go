package main

import (
	"encoding/json"
	"golang.org/x/sys/windows"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
	"time"
)

// Returns true if game is installed, false if not.
// Decided by presence of key version in gameInfo file.
func isInstalled() bool {
	// Check if gameInfo file exists
	_, err := os.Stat(gameInfo)
	if err != nil {
		// Either we can't open file, or it does not exist.
		slog.Info("Can't find gameInfo file.")
		return false
	}

	gameFile, err := os.Open(gameInfo)
	if err != nil {
		// Either we can't open file, or it does not exist.
		slog.Info("Can't open gameInfo file.")
		return false
	}
	defer gameFile.Close()

	gameData, err := io.ReadAll(gameFile)
	if err != nil {
		slog.Info("Can't read gameInfo file.")
	}

	// We were able to open up our file, but can we actually read it?
	var game Game
	err = json.Unmarshal(gameData, &game)
	if err != nil {
		// Can't read it.
		slog.Info("Can't read gameInfo file.")
		return false
	}
	// We can read it -> check if there is a key.
	if game.Version == "" {
		slog.Info("gameInfo file does not contain key version.")
		return false
	}

	// By matter of elimination, we should be installed.
	slog.Info("Game is installed.")
	return true
}

// Returns Latest struct with latest version of game.
func fetchLatestVersion() (Latest, error) {
	// Fetch the latest.json file from the file server.
	jsonUrl := fileServer + "latest.json"

	jsonClient := http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("GET", jsonUrl, nil)
	if err != nil {
		slog.Error("Error creating request:", "request", err)
		return Latest{}, err
	}

	resp, err := jsonClient.Do(req)
	if err != nil {
		slog.Error("Error fetching latest.json:", "json", err)
		return Latest{}, err
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var latest Latest
	jsonErr := json.Unmarshal(body, &latest)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	slog.Info("Latest version information fetched:", "version", latest)
	return latest, nil
}

// Returns Latest struct with latest version of game.
func fetchUpdates() (Updates, error) {
	// Fetch the updates.json file from the file server.
	jsonUrl := fileServer + "updates.json"

	jsonClient := http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("GET", jsonUrl, nil)
	if err != nil {
		slog.Error("Error creating request:", "request", err)
		return Updates{}, err
	}

	resp, err := jsonClient.Do(req)
	if err != nil {
		slog.Error("Error fetching latest.json:", "json", err)
		return Updates{}, err
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var updates Updates
	jsonErr := json.Unmarshal(body, &updates)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	slog.Info("Latest updates fetched:", "version", updates)
	return updates, nil
}

// Updates version in gameInfo file.
func updateVersion(newVersion string) error {
	// Read file
	gameFile, err := os.Open(gameInfo)
	if err != nil {
		slog.Error("Can't open gameInfo file:", "file", err)
		return err
	}
	defer gameFile.Close()

	// Unmarshal
	var game Game
	gameFileR, _ := io.ReadAll(gameFile)
	if err = json.Unmarshal(gameFileR, &game); err != nil {
		return err
	}

	// Modify the version
	game.Version = newVersion

	// Marshal back to JSON with indentation for nice formatting
	updatedJSON, err := json.MarshalIndent(game, "", "  ")
	if err != nil {
		slog.Error("Can't marshal gameInfo file:", "file", err)
	}

	// Write back to file
	err = os.WriteFile(gameInfo, updatedJSON, 0644)
	if err != nil {
		return err
	}
	return nil
}

func downloadLatest(latest Latest, keepDownload bool) bool {
	// Download the latest version.
	err := downloadFile(latest.File, "")
	if err != nil {
		slog.Error("Error downloading file:", "file", err)
		return false
	}

	// Verify the checksum.
	if !verifyChecksum(latest.File, latest.Checksum) {
		slog.Error("Checksum mismatch:", "expected", latest.Checksum, "got", "")
		slog.Error("Please delete ZIP-file to try download again.")
		return false
	}

	// We downloaded the file, now let us extract it.
	err = extractZip(latest.File)

	if !keepDownload {
		// Remove the file.
		err := os.Remove(latest.File)
		if err != nil {
			slog.Error("Error removing file:", "file", err)
		}
	} else {
		slog.Info("keep_downloads is set to true, not removing downloaded file.", "keep_downloads", keepDownload)
	}

	// Finally, let's create a gameInfo file.
	game := Game{
		Path:    "",
		Version: latest.Version,
	}
	gameJSON, err := json.MarshalIndent(game, "", "  ")
	if err != nil {
		slog.Error("Error creating gameInfo file:", "file", err)
		return false
	}

	err = os.WriteFile(gameInfo, gameJSON, 0644)
	if err != nil {
		slog.Error("Error creating gameInfo file:", "file", err)
		return false
	}
	return true

}

// Returns true if game is installed, false if not.
// Decided by presence of key version in gameInfo file.
func getGameInfo() (Game, error) {
	// Check if gameInfo file exists
	_, err := os.Stat(gameInfo)
	if err != nil {
		// Either we can't open file, or it does not exist.
		slog.Info("Can't find gameInfo file.")
		return Game{}, err
	}

	gameFile, err := os.Open(gameInfo)
	if err != nil {
		// Either we can't open file, or it does not exist.
		slog.Info("Can't open gameInfo file.")
		return Game{}, err
	}
	defer gameFile.Close()

	gameData, err := io.ReadAll(gameFile)
	if err != nil {
		slog.Info("Can't read gameInfo file.")
		return Game{}, err
	}

	// We were able to open up our file, but can we actually read it?
	var game Game
	err = json.Unmarshal(gameData, &game)
	if err != nil {
		// Can't read it.
		slog.Info("Can't read gameInfo file.")
		return Game{}, err
	}

	return game, nil
}

func updateLoop(updates Updates, game Game, keepDownload bool) {
	for {
		foundUpdate := false
		for _, update := range updates {
			if update.Update.AppliesTo == game.Version {
				// Update found! Let's download this update.
				downloadFile(update.Update.File, "")
				extractZip(update.Update.File)
				slog.Info("Update downloaded and extracted.", "version", update.Update.Version)
				if !keepDownload {
					err := os.Remove(update.Update.File)
					if err != nil {
						slog.Error("Error removing file:", "file", err)
					}
				} else {
					slog.Info("keep_downloads is set to true, not removing downloaded file.", "keep_downloads", keepDownload)
				}
				updateVersion(update.Update.Version)

				// Update our local game info after applying the update
				var err error
				game, err = getGameInfo()
				if err != nil {
					slog.Error("Failed to get updated game info:", "error", err)
					return
				}

				foundUpdate = true
				break // Start over with the new version
			}
		}

		if !foundUpdate {
			break // No more updates apply to our version
		}
	}
}

func launchGame() error {
	// Get the game exe file.
	exePath := filepath.Join("bin", "client", "Sb_client.exe")

	// Convert to absolute path.
	absPath, err := filepath.Abs(exePath)
	if err != nil {
		slog.Error("Failed to get absolute path:", "error", err)
		return err
	}

	slog.Info("Launching game:", "path", absPath)

	exePath = absPath

	cmd := exec.Command(exePath)
	cmd.Dir = filepath.Dir(exePath)

	if runtime.GOOS == "windows" {
		verb := "runas"
		verbPtr, _ := syscall.UTF16PtrFromString(verb)
		exePath, _ := syscall.UTF16PtrFromString(exePath)
		argPtr, _ := syscall.UTF16PtrFromString("")
		cwdPtr, _ := syscall.UTF16PtrFromString("")
		err := windows.ShellExecute(0, verbPtr, exePath, argPtr, cwdPtr, 1)
		if err != nil {
			slog.Error("Failed to launch game:", "error", err)
			return err
		}
	} else {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow:    false,
			CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
			// Request elevation

		}
		err = cmd.Start()
		if err != nil {
			slog.Error("Failed to launch game:", "error", err)
			return err
		}
	}

	return nil
}
