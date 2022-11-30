package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"snippetbox/pkg/models/mysql"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"

	_ "github.com/go-sql-driver/mysql"
)

// An application struct
type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	snippets       *mysql.SnippetModel
	templateCache  map[string]*template.Template
	sessionManager *scs.SessionManager
}

func main() {
	// flag will be stored in the addr variable at runtime
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "admin:admin@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	// logger for writing information message
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	// logger for writing error messages, use stderr as destination
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	// The connection pool is closed before main() function exits
	defer db.Close()

	templateCache, err := NewTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	// Initialise a new instance of our applicaion struct, containing the deps
	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		snippets:       &mysql.SnippetModel{DB: db},
		templateCache:  templateCache,
		sessionManager: sessionManager,
	}

	// Initialize a new http.Server
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s 🚀", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

// Returns a sql.DB connection pool
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
