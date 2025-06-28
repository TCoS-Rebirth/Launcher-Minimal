package main

import (
	"fmt"
	"log/slog"
)

// todo's/info:
// -> wanneer er een game.json file is, die gebruiken. Anders zelf eentje maken.
// sowieso die gebruiken he, waarom zoude in godsnaam elke keer een full download doen?

var gameInfo string = "game.json"
var fileServer string = "https://files.spellborn.org/"

func main() {

	// Check if the game is installed or not.
	installed := isInstalled()

	if !installed {
		// Perform a clean installation.
		// Step one: get the latest version.
		latest, err := fetchLatestVersion()
		if err != nil {
			slog.Error("Error fetching latest version:", "version", err)
		}

		// Step two: pass the latest version through.
		downloadLatest(latest)

		// The game should now be installed and ready for further steps.
	}

	// Now, lets parse gameInfo file so we can see what version is installed right now.
	gameInfo, err := getGameInfo()
	if err != nil {
		slog.Error("Error fetching gameInfo:", "gameInfo", err)
	}

	// Fetch update information.
	fetchedUpdates, err := fetchUpdates()
	if err != nil {
		slog.Error("Error fetching updates:", "updates", err)
	}
	fmt.Println(fetchedUpdates)

	// TODO: Make sure the downloaded file is deleted.
	// TODO: Launch the game :-)

	updateLoop(fetchedUpdates, gameInfo)

	fmt.Println("Done.")

}
