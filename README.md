# Songlink CLI

[![Github All Releases](https://img.shields.io/github/downloads/marcusziade/songlink-cli/total.svg)](https://github.com/marcusziade/songlink-cli/releases)

A Go program that retrieves Songlink and Spotify links for a given URL using the Songlink API. The output is designed to be shared as is, allowing the receiver to both use Songlink and listen to the song preview using Spotify's embed feature.

## Features

-   Retrieves Songlink and Spotify links for a given song or album URL
-   Supports command line arguments for customizing the output format
-   Automatically copies the output to the clipboard for easy sharing
-   Includes a loading indicator to provide visual feedback during the retrieval process
-   Thoroughly tested with unit tests to ensure reliability and correctness

## Installation

### macOS

#### Homebrew

```
brew tap marcusziade/songlink-cli
brew install songlink-cli
```

#### Build

1. Clone the repository: `git clone https://github.com/marcusziade/songlink-cli.git`
2. Navigate to the project directory: `cd songlink-cli`
3. Install dependencies: `go mod download`
4. Build the executable: `go build -o songlink .`
5. Run the program: `./songlink`

### Download and Run

Go to [Releases](https://github.com/marcusziade/songlink-cli/releases) and download the appropriate version for your operating system (Linux, macOS, Windows).

## Usage

1. Copy the URL of the song or album you want to retrieve links for.
2. Run the program using one of the following commands:
    - `./songlink`: Retrieves only the Songlink URL
    - `./songlink -x`: Retrieves the Songlink URL without surrounding `<>`. For Twitter
    - `./songlink -d`: Retrieves the Songlink URL surrounded by `<>` and the Spotify URL. For Discord.
    - `./songlink -s`: Retrieves only the Spotify URL
3. The program will automatically retrieve the Songlink and/or Spotify link for the song or album and copy it to your clipboard.

## Examples

Here are a few examples of how to use the Songlink CLI:

-   Retrieve only the Songlink URL:

    ```
    ./songlink
    ```

-   Retrieve the Songlink URL without surrounding `<>` + Spotify embed:

    ```
    ./songlink -x
    ```

-   Retrieve the Songlink URL surrounded by `<>` + Spotify embed:

    ```
    ./songlink -d
    ```

-   Retrieve only the Spotify URL:
    ```
    ./songlink -s
    ```

## Contributions

I welcome contributions to the Songlink CLI project! If you have any ideas, suggestions, or bug reports, please don't hesitate to open an issue or submit a pull request. To contribute:

1. Fork the repository
2. Create a new branch for your feature or bug fix
3. Make your changes and commit them with descriptive commit messages
4. Push your changes to your forked repository
5. Submit a pull request to the main repository

I appreciate your help in making this project better!

## License

This project is licensed under the [MIT License](LICENSE).

---

I hope you find this tool useful! If you have any questions or need further assistance, please let me know.
