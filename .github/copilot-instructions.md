# Copilot Instructions for DBaker

## Project Overview
- **DBaker** is a CLI tool for generating and inserting fake data into PostgreSQL databases.
- The workflow is split into two main steps:
  1. **Introspect**: Reads the schema from a live database and writes a JSON recipe (intermediate representation) describing tables/columns.
  2. **Generate**: Reads the recipe, generates fake data, and inserts it into the database.
- The codebase is modular, with clear separation between CLI commands, data generation, database adapters, and configuration.

## Key Components
- `cmd/dbaker/`: CLI entrypoints using [cobra](https://github.com/spf13/cobra). Each command (e.g., `generate`, `introspect`, `dryrun`) is defined in its own file.
- `pkg/action/`: Implements high-level actions (e.g., `Generate`, `Introspect`). Actions encapsulate workflows and are invoked by CLI commands.
- `pkg/adapter/`: Database adapters (currently PostgreSQL). Handles DB connections, schema introspection, and row insertion.
- `pkg/model/`: Data structures for tables, columns, and types. Used throughout the codebase for schema and data representation.
- `pkg/generator/`: Logic for generating fake data values for each column type.
- `pkg/config/`: Configuration structs and CLI flag bindings.
- `test/`: Integration test resources (e.g., docker-compose for Postgres, SQL init scripts).

## Developer Workflows
- **Build:** Use `make` or `go build ./cmd/dbaker` to build the CLI.
- **Run:** Use the CLI with commands like `./dbaker introspect ...` and `./dbaker generate ...`.
- **Test:** (Planned) Integration tests will use Dockerized Postgres. See `test/` for setup.
- **Debug:** Use `fmt.Printf` for tracing queries and generated values (see `WriteRow`).

## Project-Specific Patterns
- **Intermediate Representation:** All data generation is based on a JSON recipe file produced by introspection. This decouples schema discovery from data generation.
- **Table Naming:** Table references use the format `schema.table`. Parsing is handled by `splitTableName` in `pkg/action/introspect.go`.
- **Column/Placeholder Generation:** Functions like `inferColNames` and `inferPgValPlaceholders` in `pkg/adapter/postgres.go` generate SQL fragments dynamically.
- **Extensibility:** New commands follow the pattern in `cmd/dbaker/` and should bind flags to `config.Config`.

## External Dependencies
- [cobra](https://github.com/spf13/cobra) for CLI
- [gofakeit](https://github.com/brianvoe/gofakeit) for fake data generation
- Standard Go database/sql for DB access

## Examples
- To add a new CLI command, create a new file in `cmd/dbaker/`, define a `*cobra.Command`, bind flags to `config.Config`, and add the command to the root.
- To support a new DB, add a new adapter in `pkg/adapter/` and update actions to use it.

## Conventions
- Use singular, lowercase package names.
- Keep business logic out of `main.go` and CLI files; delegate to actions and adapters.
- Use GoDoc comments for exported types and functions.
- Prefer explicit error handling and return wrapped errors with context.

---
If you are unsure about a workflow or pattern, check the corresponding package or ask for clarification.
