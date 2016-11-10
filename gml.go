package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"strings"
)

// GML from
// https://kartkatalog.geonorge.no/metadata/kartverket/grunnskoler/db4b872f-264d-434c-9574-57232f1e90d2

type FeatureCollection struct {
	XMLName        xml.Name        `xml:"FeatureCollection"`
	FeatureMembers []FeatureMember `xml:"featureMember"`
}

type FeatureMember struct {
	XMLName xml.Name `xml:"featureMember"`
	School  School   `xml:"Skole"`
}

type School struct {
	Point        Point   `xml:"posisjon>Point"`
	Persons      Persons `xml:"antall_personer>AntallPersoner"`
	CountyName   string  `xml:"kommunenavn"`
	CountNymber  int     `xml:"kommunenummer"`
	BuildingType string  `xml:"bygningstype"`
	Address      Address `xml:"postadresse > PostnummerområdeId"`
	Name         string  `xml:"skolenavn"`
}

type Point struct {
	Pos  string `xml:"pos"`
	Lat  string
	Long string
}

type Address struct {
	PostCode     int    `xml:"postnummer"`
	PostLocation string `xml:"poststed"`
}

type Persons struct {
	XMLName           xml.Name `xml:"AntallPersoner"`
	NumberOfEmployees int      `xml:"antall_ansatte"`
	NumberOfStudents  int      `xml:"antall_elever"`
}

func parseSchoolLocations() ([]School, error) {
	filename := "Offentligetjenester_0000_Norge_4258_Skoler_GML.gml"
	f, err := os.Open(filename)
	if err != nil {
		return []School{}, err
	}

	fc := FeatureCollection{}

	d := xml.NewDecoder(f)
	err = d.Decode(&fc)
	if err != nil {
		return []School{}, err
	}

	schools := fc.GetSchools()
	schools, err = formatPositions(schools)
	if err != nil {
		return []School{}, err
	}

	return schools, nil
}

// Parses the positions (alat-long) and
func formatPositions(schools []School) ([]School, error) {
	for i, school := range schools {

		pos := school.Point.Pos
		position := strings.Split(pos, " ")
		fmt.Println(school, pos)
		if len(pos) < 2 {
			return []School{}, errors.New("Error: Could not parse position:" + pos)
		}

		long := position[0]
		lat := position[1]
		schools[i].Point.Lat = lat
		schools[i].Point.Long = long
	}
	return schools, nil
}

func (fc FeatureCollection) GetSchools() []School {
	schools := []School{}
	for _, fm := range fc.FeatureMembers {
		schools = append(schools, fm.School)
	}
	return schools

}

func (s *School) Print() {
	fmt.Println(s.Name)
	fmt.Println("Location:", s.Point.Lat, s.Point.Long)
}
