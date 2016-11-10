package main

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	//_, err := getSchools()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}

	schools, err := parseSchoolLocations()
	for _, school := range schools {
		school.Print()
	}
	fmt.Println(err)
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
			fmt.Println(school)

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
