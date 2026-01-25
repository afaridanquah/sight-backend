package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main() {
	files, err := os.ReadDir("/Users/afaridanquah/code/golang/sight/exports")
	if err != nil {
		fmt.Printf("failed to read directory: %v", err)
	}
	fmt.Printf("first file %+v", files[0])

	filePath := filepath.Join("/Users/afaridanquah/code/golang/sight/exports", files[0].Name())
	file, err := os.ReadFile(filePath)
	var napp Data
	if err := json.Unmarshal(file, &napp); err != nil {
		fmt.Printf("failed to read file: %v", err)
	}

	fmt.Printf("content %v", napp)
}

type Data struct {
	Identification []Identification `json:"data"`
}

type Identification struct {
	UID                int       `json:"uid"`
	Forenames          string    `json:"forenames"`
	Surname            string    `json:"surname"`
	PrevOrMaidenName   string    `json:"prev_or_maiden_name"`
	Sex                string    `json:"sex"`
	Occupation         string    `json:"occupation"`
	MaritalStatus      string    `json:"marital_status"`
	DateOfBirth        string    `json:"date_of_birth"`
	BirthTown          string    `json:"birth_town"`
	BirthCountry       string    `json:"birth_country"`
	BirthRegion        string    `json:"birth_region"`
	BirthDistrict      string    `json:"birth_district"`
	Nationality        string    `json:"nationality"`
	Resident           string    `json:"resident"`
	SocialSecurityNo   string    `json:"social_security_no"`
	MotherMaidenName   string    `json:"mother_maiden_name"`
	MotherForename     string    `json:"mother_forename"`
	NationalID         string    `json:"national_id"`
	IssueDate          string    `json:"issue_date"`
	ExpiryDate         string    `json:"expiry_date"`
	CountryOfIssue     string    `json:"country_of_issue"`
	PlaceOfIssue       string    `json:"place_of_issue"`
	CardNumber         string    `json:"card_number"`
	HouseNumber        string    `json:"house_number"`
	StreetName         string    `json:"street_name"`
	Town               string    `json:"town"`
	City               string    `json:"city"`
	Community          string    `json:"community"`
	Country            string    `json:"country"`
	Region             string    `json:"region"`
	District           string    `json:"district"`
	PostalAddress      string    `json:"postal_address"`
	PhoneNumber1       string    `json:"phone_number_1"`
	PhoneNumber2       any       `json:"phone_number_2"`
	Email              string    `json:"email"`
	TinNumber          string    `json:"tin_number"`
	TinIssuedDate      any       `json:"tin_issued_date"`
	LastUpdate         time.Time `json:"last_update"`
	CreatedDate        time.Time `json:"created_date"`
	DigLongitude       any       `json:"dig_Longitude"`
	DigLatitude        any       `json:"dig_Latitude"`
	DigStreet          any       `json:"dig_Street"`
	DigRegion          any       `json:"dig_Region"`
	DigArea            any       `json:"dig_Area"`
	DigDistrict        any       `json:"dig_District"`
	DigPostCode        any       `json:"dig_PostCode"`
	GhanaPostAddress   string    `json:"ghana_post_address"`
	LocationUID        any       `json:"location_uid"`
	RefisicID          any       `json:"refisicID"`
	RegistrationSTATUS string    `json:"registration_STATUS"`
}
