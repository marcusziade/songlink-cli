# Songlink CLI
This project is a simple Go program that retrieves Songlink and Spotify links for a given URL using the Songlink API. The output is meant to be shared as is, so the receiver can both use Songlink and listen to the song preview using Spotify's embed feature.
## Installation

### MacOS
#### Homebrew
```
brew tap marcusziade/songlink-cli
brew install songlink-cli
```
#### Build
1. Clone the repository: `git clone https://github.com/marcusziade/songlink-cli.git`
2. Install dependencies: `go mod download`
3. Build the executable: `go build -o songlink .`
4. Run the program: `./songlink`

### Download and run
Go to [Releases](https://github.com/marcusziade/songlink-cli/releases) (Linux, macOS, Windows)

## Usage
1. Copy the URL of the song or album you want to retrieve links for
2. Run the program using the command ./songlink
3. The program will automatically retrieve the Song.link and Spotify link for the song or album and copy it to your clipboard


## Contributions
Fork and make a PR or create an issue.
