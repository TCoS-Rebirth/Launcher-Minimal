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

	// To be certain - since installation leads to no version number being mentioned - re-read the version number.
	gameInfo, err = getGameInfo()

	// Fetch update information.
	fetchedUpdates, err := fetchUpdates()
	if err != nil {
		slog.Error("Error fetching updates:", "updates", err)
	}

	// Start the update loop.
	slog.Info("Starting update loop.")
	updateLoop(fetchedUpdates, gameInfo, gameInfo.KeepDownloads)
	slog.Info("Finished update loop.")

	if !gameInfo.NoLaunch {
		slog.Info("Finished update loop, launching game.")
		if err := launchGame(); err != nil {
			slog.Error("Error launching game:", "game", err)
		}
		slog.Info("Finished launching game. Goodbye! <3", "game", "launched")
	} else {
		slog.Info("Finished update loop, but not launching game as no_launch is true. Goodbye! <3", "game", "not launched")
	}
}
