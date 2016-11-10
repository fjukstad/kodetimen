package main

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/paulmach/go.geojson"
)

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)
	r.HandleFunc("/schools", SchoolsHandler)

	fmt.Println("Server started on localhost:8000")
	err := http.ListenAndServe(":8000", r)

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
	return
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	indexTemplate := template.Must(template.ParseFiles("index.html"))
	indexTemplate.Execute(w, nil)
}
