package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"arnoldcodes.com/snippetbox/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
)

//application struct to hold the Application wide dependencies
//for the web application
type Application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	//define a new command line flag for tje MySQL DSN string
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MYSQL database")

	//define a new command-line flag wuth name 'addr', a default value of the server
	//and some short help text exaplaining what the flag controls. The value of
	//flag will be stored in the addr variable at runtime.
	//addr := flag.String("addr", ":4000", "HTTP network address")
	//Importantly, we use the flag.Parse() function to parse the command-line
	//This reads in the command-line flag value and assigns it to the addr variable
	//otherwise it will alwyas contain the default value of ":4000". If any error is
	//encountered during parsing the application will be terminated
	//	flag.Parse()
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//to keep the main() function tidy we put the code for creating a contection
	//pool into the seperate openDB() function below. We pass openDB() the DSN
	//from the command line flag
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	//we also defer a call to db.Close(), so that the connection pool is close
	//before the main() function exits
	defer db.Close()

	// Initialize a new template cache...
	templateCache, err := newTemplateCache("ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}
	//new instance of application containing the depenceis of Application
	app := &Application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
	}
	//initialize a new http.Server struct.
	srv := &http.Server{
		Addr:     GetPort(),
		ErrorLog: errorLog,
		Handler:  app.routes(), //call the new app.routes() method
	}

	//Write message using the two new loggers,instead of the standard logger
	infoLog.Printf("Starting server on %s", GetPort())
	err = srv.ListenAndServe()
	errorLog.Fatal(err)

	//start server
	//the value returned from the flag.String() function is a pointer to the flag
	//value, not the value itself. So we need to dereference the pointer (i.e.
	//prefix it with the * symbol) before using it.
	/*log. Printf("Starting server on %s", *addr)
	err := http. ListenAndServe(*addr, mux)
	log.Fatal(err)*/

}
func GetPort() string {
	var port = os.Getenv("PORT")
	//set default port if the is nothing in the environmental port
	if port == "" {
		port = "4000"
		fmt.Println("INFO: No port environment variable detected, defaulting to: " + port)
	}
	return ":" + port
}

//mysql database
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
