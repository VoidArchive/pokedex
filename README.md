# Pokedex CLI

<!--toc:start-->
- [Pokedex CLI](#pokedex-cli)
  - [Features](#features)
  - [Getting Started](#getting-started)
    - [Prerequisites](#prerequisites)
    - [Building the Application](#building-the-application)
    - [Running the Application](#running-the-application)
  - [Commands](#commands)
  - [Data Persistence](#data-persistence)
  - [Cache](#cache)
<!--toc:end-->

A command-line interface (CLI) application to interact with a Pokedex. You can explore Pokemon locations, catch Pokemon, manage your Pokedex and party, and even simulate battles!

## Features

- Explore different location areas to find Pokemon.
- Catch Pokemon using different types of Pokeballs.
- Inspect your caught Pokemon to see their stats, types, level, and XP.
- Manage your Pokedex (all caught Pokemon) and your active party (up to 6 Pokemon).
- Simulate battles between your Pokemon and wild opponents.
- Pokemon can gain XP and level up from battles.
- Level-up and time-based evolutions are implemented.
- Your Pokedex, party, and inventory are saved to `pokedex.json` when you exit and loaded when you start.

## Getting Started

### Prerequisites

- Go (version 1.21 or later recommended - current `go.mod` uses `go 1.24.3`)

### Building the Application

1. **Clone the repository (if you haven't already):**

    ```bash
    # git clone <repository-url>
    # cd pokedex
    ```

2. **Build the executable:**
    Open your terminal in the project's root directory (`pokedex`) and run:

    ```bash
    go build
    ```

    This will create an executable file named `pokedex` (or `pokedex.exe` on Windows) in the current directory.

### Running the Application

Once built, you can run the Pokedex CLI from the project's root directory:

```bash
./pokedex
```

You will see the `Pokedex >` prompt.

## Commands

Type `help` at the prompt to see a list of available commands and their descriptions:

- `help`: Displays a help message.
- `exit`: Exits the Pokedex (saves your game data).
- `map`: Displays the next 20 Pokemon location areas.
- `mapb`: Displays the previous 20 Pokemon location areas.
- `explore <location_area_name_or_number>`: Explore a location area for Pokemon. You can use the number from the `map` command output.
- `catch <pokemon_name> [pokeball_type]`: Attempt to catch a Pokemon (e.g., `catch pikachu greatball`). Defaults to "pokeball".
- `inspect <pokemon_name>`: View details of a caught Pokemon.
- `pokedex`: View all Pokemon in your Pokedex.
- `party`: View the Pokemon in your active party.
- `inventory`: View your items, including Pokeballs.
- `battle <your_pokemon> <opponent_pokemon>`: Simulate a battle.

## Data Persistence

Your Pokedex (all caught Pokemon), current party, and inventory are automatically saved to a file named `pokedex.json` in the root of the project directory when you exit the application using the `exit` command, or by pressing `Ctrl+C` or `Ctrl+D`. This data is loaded the next time you start the Pokedex.

The `pokedex.json` file is included in the `.gitignore` file to prevent accidental commits of personal game data.

## Cache

The application uses a cache for API responses to speed up subsequent requests for the same data and to be mindful of API rate limits. Cache entries expire after a set interval (currently 5 minutes).

Enjoy your Pokedex adventure!
