package main

import (
	"github.com/B64-Cryptzo/MotoPi/backend/API"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers for all requests
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

func main() {

	router := httprouter.New()

	_ = API.NewHALInterfaceHandler(&API.StubHALService{}, router)
	_ = API.NewNetworkInterfaceHandler(&API.StubNetworkService{}, router)
	_ = API.NewMotorcycleInterfaceHandler(&API.StubMotorcycleService{}, router)

	log.Println("Starting backend on :8080")
	log.Fatal(http.ListenAndServe(":8080", corsMiddleware(router)))
}
