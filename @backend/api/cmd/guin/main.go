package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

var start = 0
var end = 1000

var target = 17272615
var workers = 3
var collected = 0

type Response struct {
	RecordsFiltered int `json:"recordsFiltered"`
	Data            []struct {
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
	} `json:"data"`
	Start        int `json:"start"`
	Draw         int `json:"draw"`
	RecordsTotal int `json:"recordsTotal"`
}

func main() {
	// var wg sync.WaitGroup
	fmt.Println("Starting............")

	if err := fetchData(start, end); err != nil {
		fmt.Printf("fetchdata: %v", err)
	}
	// for range workers {
	// 	wg.Add(1)
	// 	if collected >= target {
	// 		return
	// 	}
	// 	go func() {
	// 		defer wg.Done()
	// 		if err := fetchData(start, end); err != nil {
	// 			panic(1)
	// 		}
	// 		fmt.Printf("url: %d", end)

	// 		start = end
	// 		end = start + 1000
	// 		collected = end
	// 	}()
	// }

	// for i := 0; i < workers; i++ {
	// 	wg.Add(1)
	// 	go func() {
	// 		defer wg.Done()
	// 		if err := fetchData(start, end); err != nil {
	// 			fmt.Printf("fetchdata: %v", err)
	// 		}
	// 	}()
	// }

}

func fetchData(start int, end int) error {

	search := "{\"value\":\"\",\"regex\":false}"
	// endpoint := fmt.Sprintf("https://guinportal.gra.gov.gh/GetAllNiaData?draw=1start=%dlength=%dsearch=%s", start, end, search)
	endpoint := "https://guinportal.gra.gov.gh/GetAllNiaData"

	pu, err := url.Parse(endpoint)
	if err != nil {
		return err
	}
	q := pu.Query()
	q.Add("draw", "1")
	q.Add("start", string(start))
	q.Add("length", string(end))
	q.Add("search", search)

	pu.RawQuery = q.Encode()

	fmt.Printf("url: %s", pu.String())
	req, err := http.NewRequest("GET", pu.String(), nil) // nil for no request body
	if err != nil {
		return err
	}

	// Add custom headers
	req.Header.Add("Cookie", "JSESSIONID=6E2F0CF8EEA34C465EDE0B1C8925FDCF; _ga=GA1.3.341035717.1750376383; _ga_ZF0KSVVM14=GS2.3.s1754539045$o3$g1$t1754539064$j41$l0$h0; _gid=GA1.3.482285717.1754539045; _gat=1")
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("User-Agent", "insomnia/11.0.2")

	// Create an HTTP client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			return
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("body: %v", body)

	var response Response

	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		return err
	}
	// filename := fmt.Sprintf("%d.json", end)
	if err := os.WriteFile("output.json", jsonData, 0644); err != nil {
		return err
	}

	return nil

}
