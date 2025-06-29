package main

import (
	"log/slog"
)

// Configuration options
var gameInfo string = "game.json"
var fileServer string = "https://files.spellborn.org/"

func main() {

	// Check if the game is installed or not.
	installed := isInstalled()

	// Get game information, if it is present.
	gameInfo, err := getGameInfo()
	if err != nil {
		slog.Error("Error fetching gameInfo:", "gameInfo", err)
	}

	if !installed {
		// Perform a clean installation.
		// Step one: get the latest version.
		latest, err := fetchLatestVersion()
		if err != nil {
			slog.Error("Error fetching latest version:", "version", err)
		}

		// Step two: pass the latest version through.
		downloadLatest(latest, gameInfo.KeepDownloads)

		// The game should now be installed and ready for further steps.
	}

	// Fetch update information.
	fetchedUpdates, err := fetchUpdates()
	if err != nil {
		slog.Error("Error fetching updates:", "updates", err)
	}

	// TODO: Make sure the downloaded file is deleted.
	// TODO: Launch the game :-)

	updateLoop(fetchedUpdates, gameInfo)

	slog.Info("Finished update loop, launching game.")

}
