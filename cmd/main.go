package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	livelead "github.com/bysidecar/livelead/pkg"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	port := lookupEnv("DB_PORT")
	portInt, err := strconv.ParseInt(port, 10, 64)
	if err != nil {
		log.Fatalf("Error parsing to string Datbase's port %s, Err %v", port, err)
	}

	database := &livelead.Database{
		Host:      lookupEnv("DB_HOST"),
		Port:      portInt,
		User:      lookupEnv("DB_USER"),
		Pass:      lookupEnv("DB_PASS"),
		Dbname:    lookupEnv("DB_NAME"),
		Charset:   "utf8",
		ParseTime: "True",
		Loc:       "Local",
	}

	if err := database.Open(); err != nil {
		log.Fatalf("error opening database connection. err %v", err)
	}
	defer database.Close()

	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("error creating table, err %v", err)
	}

	c := livelead.Client{
		Storer: database,
	}

	router := mux.NewRouter()

	router.PathPrefix("/lead/live").Handler(c.HandleFunction()).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(":4500", cors.Default().Handler(router)))
}

// lookupenv looks for an environment variable
func lookupEnv(envvar string) string {
	env, ok := os.LookupEnv(envvar)
	if !ok {
		log.Fatalf("Error: %s ENV var not found", envvar)
	}
	return env
}
