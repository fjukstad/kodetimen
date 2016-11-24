package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type RegisteredSchool struct {
	School
	Registered bool
}
type LKKSchool struct {
	Previous interface{} `json:"_previous"`
	Locality string      `json:"locality"`
	Hash     string      `json:"_hash"`
	Ts       int64       `json:"_ts"`
	ID       string      `json:"_id"`
	PostDate string      `json:"post_date"`
	Updated  int         `json:"_updated"`
	Level    string      `json:"level"`
	County   string      `json:"county"`
	Deleted  bool        `json:"_deleted"`
	Address  string      `json:"address"`
	PosLat   string      `json:"pos_lat"`
	Students interface{} `json:"students"`
	School   string      `json:"school"`
	Year     int         `json:"year"`
	PosLong  string      `json:"pos_long"`
}

func getRegisteredSchools() ([]RegisteredSchool, error) {
	registeredSchools, err := getSchools()
	if err != nil {
		fmt.Println(err)
		return []RegisteredSchool{}, err
	}

	allNorwegianSchools, err := parseSchoolLocations()
	if err != nil {
		return []RegisteredSchool{}, err
	}

	schools := []RegisteredSchool{}

	for _, school := range allNorwegianSchools {
		reg := false
		for _, registered := range registeredSchools {
			if strings.Contains(school.Name, registered.School) && registered.Locality == school.MunicipalityName {
				reg = true
				break
			}
		}
		schools = append(schools, RegisteredSchool{school, reg})
	}
	return schools, nil
}

func getSchools() ([]LKKSchool, error) {
	resp, err := http.Get("https://open.sesam.io/api/api/datasets/skoler_med_kodetimen/entities")

	if err != nil {
		return []LKKSchool{},
			errors.Wrap(err, "Could not fetch data from open.sesam.io")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []LKKSchool{}, errors.Wrap(err, "Could not read http body")
	}

	var schools []LKKSchool
	err = json.Unmarshal(body, &schools)

	if err != nil {
		return []LKKSchool{}, errors.Wrap(err, "Could not unmarshal json from open.sesam.io")
	}

	return schools, nil
}

func inSlice(str string, slice []string) bool {
	for _, a := range slice {
		if str == a {
			return true
		}
	}
	return false

}

func sortSchools(schools []RegisteredSchool) []RegisteredSchool {
	var registered, notRegistered []RegisteredSchool
	for _, school := range schools {
		if school.Registered {
			registered = append(registered, school)
		} else {
			notRegistered = append(notRegistered, school)
		}
	}
	return append(notRegistered, registered...)
}
