package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/savingoyal/systems-assignment/pkg/api/v1alpha1"
)

type Server struct {
	port int
	host string
	data *v1alpha1.KeyValueStore
}

func Serve(port int, host string, dataFilename string) {
	// Create a new server
	s := NewServer(port, host, dataFilename)
	s.Start()
}

func NewServer(port int, host string, dataFilename string) *Server {
	// Create a new server
	s := &Server{
		port: port,
		host: host,
	}

	// Load the data from the file
	data, err := v1alpha1.LoadKeyValueStore(dataFilename)
	if err != nil {
		log.Fatalf("Error loading data: %v", err)
	}
	s.data = data

	return s
}

func (s *Server) Start() {
	// Start the server
	fmt.Printf("Starting server on %s:%d\n", s.host, s.port)

	router := mux.NewRouter().StrictSlash(true)
	router.Use(loggingMiddleware)

	router.HandleFunc("/", homePage).Methods("GET")
	router.HandleFunc("/api", homePage).Methods("GET")
	router.HandleFunc("/api/lookup", homePage).Methods("GET")
	router.HandleFunc("/api/lookup/{id}", s.apiLookup).Methods("GET", "HEAD")
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", s.host, s.port), router))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Printf("%s %s %s", r.Method, r.RequestURI, r.Proto)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func (s *Server) Get(key string) (*v1alpha1.KeyValue, error) {
	// Get the value for the key
	value, err := s.data.Get(key)
	return value, err
}

//nolint:errcheck // ignore unchecked errors that cannot be returned
func (s *Server) apiLookup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	// Get the value for the key
	kv, err := s.Get(key)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method == "HEAD" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Write the value to the response
	if acceptContains(r.Header.Get("Accept"), "application/json") {
		result, err := json.Marshal(kv)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(result)
	} else {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(kv.Value))
	}
}

//nolint:errcheck // ignore unchecked errors that cannot be returned
func homePage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if acceptContains(r.Header.Get("Accept"), "text/html") {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<html><body>"))
		w.Write([]byte("<h2>Welcome to the KeyValueStore API!<h2>"))
		w.Write([]byte("<h3>valid endpoints:</h3>"))
		w.Write([]byte("<ul><li><div>/api/lookup/{ID}</div>"))
		w.Write([]byte("<div style=\"margin-left: 3em\">Lookup the value for ID</div></li></ul>"))
		w.Write([]byte("</body></html>"))
	} else {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Welcome to the KeyValueStore API!\n\n"))
		w.Write([]byte("valid endpoints:\n"))
		w.Write([]byte("\t/api/lookup/{ID}\n"))
		w.Write([]byte("\t\tLookup the value for ID\n"))
	}
}

func acceptContains(accept string, mimeType string) bool {
	for _, a := range strings.Split(accept, ",") {
		if a == mimeType {
			return true
		}
	}
	return false
}
