package countryapp

type Country struct {
	Name   string `json:"name"`
	Alpha2 string `json:"cca2"`
	Alpha3 string `json:"cca3"`
	Emoji  string `json:"emoji"`
}
