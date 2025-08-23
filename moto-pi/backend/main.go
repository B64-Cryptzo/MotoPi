package main

import (
	"log"
	"net/http"

	"github.com/B64-Cryptzo/MotoPi/backend/API"
	"github.com/B64-Cryptzo/MotoPi/backend/Firmware/hal/gps"
	"github.com/B64-Cryptzo/MotoPi/backend/Firmware/hal/rfid"
	"github.com/julienschmidt/httprouter"
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

	scanner := &rfid.RFIDScanner{}
	if err := scanner.Init(); err != nil {
		panic(err)
	}

	gps := gps.NewGPS("/dev/ttyUSB0", 9600)
	if err := gps.Init(); err != nil {
		panic(err)
	}

	defer gps.Close()
	defer scanner.Close()

	router := httprouter.New()

	_ = API.NewHALInterfaceHandler(&API.LiveHALService{RFIDScanner: scanner, GPS: gps}, router)
	_ = API.NewNetworkInterfaceHandler(&API.StubNetworkService{}, router)
	_ = API.NewMotorcycleInterfaceHandler(&API.StubMotorcycleService{}, router)

	log.Println("Starting backend on :8080")
	log.Fatal(http.ListenAndServe(":8080", corsMiddleware(router)))
}
