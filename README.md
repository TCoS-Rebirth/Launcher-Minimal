# Minimal launcher for The Chronicles of Spellborn - Reborn
This is an alternative **minimal** and **multi-platform** version of the launcher used for [The Chronicles of Spellborn Reborn](https://spellborn.org).
Unlike the main launcher, there is no built-in browser. Due to this it is extremely small (about 6-7MB) and requires no runtime or Edge WebView components installed.

## Configuration options
The launcher will create a file called `game.json` if none is present. This can currently contain the following values:
* `path` (string), currently un-used. Meant to be used to install the game to another installation folder.
* `version` (string), filled in by the launcher when a file has been succesfully extracted. Used to determine the current game version.
* `keep_downloads` (boolean, default: false), stops the launcher from removing the downloaded files after they are extracted.
* `no_launch` (boolean, default: false), stops the launcher from launching the game executeable after the update loop is finished.

## Why would you use this launcher instead?
- It is way more lightweight (6-7MB).
- It is fully open source, unlike the other launcher (it is my first and only C# GUI project, and it is messy).
- It supports **resuming of downloads** unlike the main launcher.
- It is shiny and written in Golang <3
- Due to it being written in Golang, it requires no runtimes or extra files. A single binary!
- It does not use the registry and instead uses a file called (by default) game.json
- It runs on way more machines compared to the other launcher which requires Windows 10+.
- It runs on Mac OS and Linux (you'll still need Wine or an emulator to run Spellborn itself).

## Building
In order to trigger a new build, create a new tag and Goreleaser will create binaries for that tag.

## Suggestions? Improvements?
Feel free to let me know here or on Discord.

Have fun and I'll see you in game!
