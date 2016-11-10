package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type RegisteredSchool struct {
	School
	Registered bool
}

func getRegisteredSchools() ([]RegisteredSchool, error) {
	registeredSchools, err := getSchools()
	if err != nil {
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
			if registered == school.Name {
				reg = true
				break
			}
		}
		schools = append(schools, RegisteredSchool{school, reg})
	}
	return schools, nil
}

func getSchools() ([]string, error) {
	url := "http://kidsakoder.no/kodetimen/pameldte-skoler-2016/"
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return []string{}, err
	}

	var schools []string

	// Read list of school from Kidsakoder.no and append to list of schools if
	// not already in list (some schools are registered more than once).
	doc.Find(".kodetimen-item").Each(func(i int, s *goquery.Selection) {
		school := s.Find("span").Text()
		school = strings.TrimSpace(school)
		school = strings.Split(school, ",")[0]
		if !inSlice(school, schools) {
			schools = append(schools, school)
		}
	})
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
