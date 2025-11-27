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
	// var start = 0
	// var end = 1000
	fmt.Println("Starting............")

	for i := 0; i < 17272615; i += 20000 {
		if err := fetchData(i); err != nil {
			fmt.Printf("err: %v\n", err)
		}
		time.Sleep(time.Second * 3)
	}

	// var wg sync.WaitGroup
	// for i := 0; i < 17272615; i += 100 {
	// 	wg.Add(1)
	// 	time.Sleep(time.Second * 5)

	// 	go func(x int) {
	// 		if err := fetchData(x); err != nil {
	// 			fmt.Printf("err: %v\n", err)
	// 		}
	// 		wg.Done()
	// 	}(i)
	// }

	// wg.Wait()

	fmt.Println("All workers have finished.")

	// if err := fetchData(start, end); err != nil {
	// 	fmt.Printf("fetchdata: %v", err)
	// }

}

func fetchData(start int) error {
	endpoint := "https://guinportal.gra.gov.gh/GetAllNiaData?draw=4&columns[0][data]=forenames&columns[0][name]=&columns[0][searchable]=true&columns[0][orderable]=true&columns[0][search][value]=&columns[0][search][regex]=false&columns[1][data]=date_of_birth&columns[1][name]=&columns[1][searchable]=true&columns[1][orderable]=true&columns[1][search][value]=&columns[1][search][regex]=false&columns[2][data]=national_id&columns[2][name]=&columns[2][searchable]=true&columns[2][orderable]=true&columns[2][search][value]=&columns[2][search][regex]=false&columns[3][data]=tin_number&columns[3][name]=&columns[3][searchable]=true&columns[3][orderable]=true&columns[3][search][value]=&columns[3][search][regex]=false&columns[4][data]=phone_number_1&columns[4][name]=&columns[4][searchable]=true&columns[4][orderable]=true&columns[4][search][value]=&columns[4][search][regex]=false&columns[5][data]=Status&columns[5][name]=&columns[5][searchable]=true&columns[5][orderable]=true&columns[5][search][value]=&columns[5][search][regex]=false"

	pu, err := url.Parse(endpoint)
	if err != nil {
		return err
	}
	q := pu.Query()
	q.Add("draw", "1")
	q.Add("start", fmt.Sprintf("%d", start))
	q.Add("length", "20000")
	q.Add("search", "{'value':'','regex':false}")

	pu.RawQuery = q.Encode()

	fmt.Printf("url: %s\n", pu.String())

	req, err := http.NewRequest("GET", pu.String(), nil) // nil for no request body
	if err != nil {
		return err
	}

	// Add custom headers
	req.Header.Add("Cookie", "_ga=GA1.3.341035717.1750376383; _ga_ZF0KSVVM14=GS2.3.s1762721238$o20$g1$t1762721263$j35$l0$h0; _ga_JXGVPSYQ59=GS2.1.s1754608648$o1$g1$t1754608690$j18$l0$h0; JSESSIONID=2E04AC6243E3D1633465A1B2EAAE3AB8; _gid=GA1.3.379434554.1762700094; _gat=1")
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	req.Header.Add("Accept", "*/*")
	// req.Header.Add("User-Agent", "insomnia/11.0.2")

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

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		return err
	}
	filename := fmt.Sprintf("exports/%d.json", start)
	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		return err
	}

	return nil
}
