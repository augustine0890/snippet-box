# Web Development
- A web application called Snippetbox. Like Pastebin or GitHub’s Gists.
## Project Structure and Organization
```
.
├── README.md
├── cmd
│   └── web
│       └── main.go
├── docker-compose.yaml
├── example.env
├── go.mod
├── pkg
└── ui
    ├── html
    └── static
```

- The `cmd` directory will contain the application specific code for the executable apps in project. The web application which will live under the `cmd/web` directory.
- The `pkg` directory will contain the non-application code used in the project. It holds potentially reusable code like validation helpers and the SQL database models for the project.
- The `ui` directory will contain the user-interface assets used by the web app.
### Logging
- Prefix informational messages with `"INFO"`
- Prefix error messages with `"ERROR"`
- Logging to a File
  ```go
  f, err := os.OpenFile("/tmp/info.log", os.O_RDWR|os.O_CREATE, 0666) if err != nil {
  log.Fatal(err) }
  defer f.Close()
  infoLog := log.New(f, "INFO\t", log.Ldate|log.Ltime)
  ```

## How to run the app
- `go run ./cmd/web`
- Using the `-addr` flag
  - `go run ./cmd/web -addr=":3000"`
- Automated help
  - `go run ./cmd/web -help`

- Pre-Existing Variables
  ```go
  type Config struct {
    Addr      string
    StaticDir string
  }

  cfg := new(Config)
  flag.StringVar(&cfg.Addr, "addr", ":4000", "HTTP network  address")
  flag.StringVar(&cfg.StaticDir, "static-dir", "./ui/static", "Path to static assets")
  flag.Parse()
  ```

- All incoming HTTP requests are served in their own goroutine.
- Race conditions are accessing the shared resources from your handlers.

## Database
### Working with Transactions
- Transactions are also super-useful if you want to execute multiple transactions at single atomic action.
- `tx.Rollback()` method in the event of any errors, the transaction ensures that either:
  - All statements are executed successfully or
  - No statements are executed an the database remains unchanged

## Panic Recovery
- Panic Server will log a stack trace to the server error log, unwind the stack for the affected goroutine and close the underlying HTTP connection.
- Do not terminate the application --> any panic will not bring down the server.