# Podcast Filter

This project is a Go application designed to filter podcast feeds based on specified criteria. It utilizes GitHub Actions to automate the filtering process and publishes the filtered feed on GitHub Pages.

## Project Structure

```
podcast-filter
├── src
│   ├── main.go          # Entry point of the application
│   └── filter
│       └── filter.go    # Contains filtering logic
├── .github
│   └── workflows
│       └── filter.yml    # GitHub Actions workflow for filtering
├── .gitignore           # Files and directories to ignore by Git
├── go.mod               # Module dependencies
├── go.sum               # Checksums for module dependencies
└── README.md            # Project documentation
```

## Setup Instructions

1. **Clone the repository:**
   ```
   git clone https://github.com/yourusername/podcast-filter.git
   cd podcast-filter
   ```

2. **Install Go:**
   Ensure you have Go installed on your machine. You can download it from [golang.org](https://golang.org/dl/).

3. **Install dependencies:**
   Run the following command to install the necessary dependencies:
   ```
   go mod tidy
   ```

4. **Run the application:**
   You can run the application using:
   ```
   go run src/main.go
   ```

## Usage

The application will read the podcast feed, apply the filtering criteria defined in `src/filter/filter.go`, and output the filtered feed.

## GitHub Actions

The filtering process is automated using GitHub Actions. The workflow defined in `.github/workflows/filter.yml` will trigger the filtering function whenever changes are pushed to the repository.

## Publishing

The filtered podcast feed will be published on GitHub Pages. Ensure that the necessary configurations are set up in your repository settings to enable GitHub Pages.