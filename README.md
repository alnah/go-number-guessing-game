# Go Number Guessing Game

This repository hosts a simple number guessing CLI game written in Go. While the functionality is straightforward—guess a number between 1 and 100—its main purpose is to serve as a toy app for practicing:

- **Testing** (unit and integration tests)
- **Error wrapping** (custom error types and context)
- **Package structuring** (separating responsibilities under `internal/`)
- **Project layout** (clear organization with `cmd`, `configs`, `Makefile`, etc.)

The idea is taken from [roadmap.sh](https://roadmap.sh/projects/number-guessing-game).

---

## Getting Started

Clone the repository.

```bash
git clone https://github.com/alnah/go-number-guessing-game.git
cd go-number-guessing-game
```

Install dependencies. Ensure you have Go installed (1.18+). Dependencies are tracked in go.mod and go.sum. You can verify or download them with:

```bash
go mod download
```

Build the binary. Use the provided Makefile:

```bash
make build
```

Run the game. After building, run the binary:

```bash
./number-guessing
```

You will be prompted to enter your name, difficulty level, and guesses.

Optionally, run the tests.

Run all tests (unit + integration):

```bash
make test
```

## Key Highlights

Project Structure

- `cmd/main.go`: The entry point for the application.
- `internal/`: Contains feature-specific sub-packages:
- `cli`: Handles user input abstraction and display utilities.
- `config`: Loads YAML configs using the Viper library.
- `game`: Core logic (turns, validation, outcomes).
- `parser`: Validates and parses user inputs.
- `service`: Orchestrates gameplay flow and integrates other packages.
- `store`: Persists and retrieves top scores from a JSON file.
- `timer`: Tracks elapsed time in a session.
- `configs/`: Stores YAML config files for the game.
- `makefile`: Basic commands for build and test automation.

Testing

- Focused Unit Tests: Each package has dedicated unit tests.
- Integration Tests: Packages like `config_test`, `service_test`, and `store_test` include integration tests for interactions between components.
- Mocks and Stubs: Simulate user input and data stores for predictable tests.

Error Wrapping

- Custom Error Types: Examples include `ReadConfigError`, `LevelError`, `TurnsLengthError`, etc., for granular test assertions and improved debugging.
- Error Context: Errors are wrapped or constructed with messages to clarify failure points.

Package Structuring

- Encapsulation: The internal directory encapsulates code not meant for external consumption, aligning with Go project structure conventions.
- Separation of Concerns: Each sub-package handles a distinct concern, making the codebase maintainable and testable.

Project Layout

- Clean Organization: Uses `cmd/` for the entry point, `internal/` for private logic, and `configs/` for configuration.
- Scalability: This layout supports adding more commands, configurations, or integrations seamlessly.
