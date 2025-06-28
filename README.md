# Minimal launcher for The Chronicles of Spellborn - Reborn
This is an alternative **minimal** and **multi-platform** version of the launcher used for [The Chronicles of Spellborn Reborn](https://spellborn.org).
Unlike the main launcher, there is no built-in browser. Due to this it is extremely small (about 6-7MB) and requires no runtime or Edge WebView components installed.

## Why would you use this launcher instead?
- It is way more lightweight (6-7MB).
- It is fully open source, unlike the other launcher (it is my first and only C# GUI project, and it is messy).
- It supports **resuming of downloads** unlike the main launcher.
- It is shiny and written in Golang <3
- It does not use the registry and instead uses a file called (by default) game.json
- It runs on way more machines compared to the other launcher which requires Windows 10+.
- It runs on Mac OS and Linux (you'll still need Wine or an emulator to run Spellborn itself).

## Building
In order to trigger a new build, create a new tag and Goreleaser will create binaries for that tag.

## Suggestions? Improvements?
Feel free to let me know here or on Discord.

Have fun and I'll see you in game!
