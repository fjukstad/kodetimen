package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/paulmach/go.geojson"
)

func main() {

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))
	mux.HandleFunc("/schools", SchoolsHandler)

	port := os.Getenv("PORT")

	if port == "" {
		fmt.Println("$PORT must be set")
	}

	fmt.Println("Server started on", port)
	err := http.ListenAndServe(":"+port, mux)

	if err != nil {
		fmt.Println(err)
		return
	}
}

func SchoolsHandler(w http.ResponseWriter, r *http.Request) {

	schools, err := getRegisteredSchools()
	if err != nil {
		w.Write([]byte("Could not get registered schools"))
		return
	}

	fc := geojson.NewFeatureCollection()

	for _, school := range schools {
		geom := geojson.NewPointGeometry([]float64{school.Point.Long, school.Point.Lat})
		f := geojson.NewFeature(geom)
		f.SetProperty("name", school.Name)
		f.SetProperty("email", school.Owner.Email)
		f.SetProperty("phonenumber", school.Owner.PhoneNumber)
		f.SetProperty("registered", school.Registered)
		fc = fc.AddFeature(f)
	}

	b, err := fc.MarshalJSON()

	if err != nil {
		fmt.Println(err)
		w.Write([]byte("Could not marshal geojson"))
		return
	}

	w.Write(b)
}
