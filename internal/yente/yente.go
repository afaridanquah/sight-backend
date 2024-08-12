package yente

const (
	YENTE_URL = "http://159.89.178.0:8000/match/default?algorithm=logic-v1"
)

type YenteResponse struct {
	Responses struct {
		Match struct {
			Status  int `json:"status"`
			Results []struct {
				ID         string `json:"id"`
				Caption    string `json:"caption"`
				Schema     string `json:"schema"`
				Properties struct {
					Position      []string `json:"position"`
					Alias         []string `json:"alias"`
					Name          []string `json:"name"`
					Notes         []string `json:"notes"`
					Gender        []string `json:"gender"`
					AddressEntity []string `json:"addressEntity"`
					BirthDate     []string `json:"birthDate"`
					BirthPlace    []string `json:"birthPlace"`
					Topics        []string `json:"topics"`
					LastName      []string `json:"lastName"`
					FirstName     []string `json:"firstName"`
					Address       []string `json:"address"`
					Education     []string `json:"education"`
					CreatedAt     []string `json:"createdAt"`
					InnCode       []string `json:"innCode"`
					SecondName    []string `json:"secondName"`
					Country       []string `json:"country"`
					Nationality   []string `json:"nationality"`
					ModifiedAt    []string `json:"modifiedAt"`
					Religion      []string `json:"religion"`
					SourceURL     []string `json:"sourceUrl"`
					WeakAlias     []string `json:"weakAlias"`
					FatherName    []string `json:"fatherName"`
					WikidataID    []string `json:"wikidataId"`
					Website       []string `json:"website"`
					Ethnicity     []string `json:"ethnicity"`
					Keywords      []string `json:"keywords"`
					Status        []string `json:"status"`
					BirthCountry  []string `json:"birthCountry"`
				} `json:"properties"`
				Datasets   []string `json:"datasets"`
				Referents  []string `json:"referents"`
				Target     bool     `json:"target"`
				FirstSeen  string   `json:"first_seen"`
				LastSeen   string   `json:"last_seen"`
				LastChange string   `json:"last_change"`
				Score      int      `json:"score"`
				Features   struct {
					NameLiteralMatch        int     `json:"name_literal_match"`
					PersonNameJaroWinkler   int     `json:"person_name_jaro_winkler"`
					PersonNamePhoneticMatch float64 `json:"person_name_phonetic_match"`
					NameSoundexMatch        float64 `json:"name_soundex_match"`
				} `json:"features"`
				Match bool `json:"match"`
			} `json:"results"`
			Total struct {
				Value    int    `json:"value"`
				Relation string `json:"relation"`
			} `json:"total"`
			Query struct {
				ID         interface{} `json:"id"`
				Schema     string      `json:"schema"`
				Properties struct {
					Name        []string `json:"name"`
					Nationality []string `json:"nationality"`
				} `json:"properties"`
			} `json:"query"`
		} `json:"match"`
	} `json:"responses"`
	Matcher struct {
		NameLiteralMatch struct {
			Description string `json:"description"`
			Coefficient int    `json:"coefficient"`
			URL         string `json:"url"`
		} `json:"name_literal_match"`
		PersonNameJaroWinkler struct {
			Description string  `json:"description"`
			Coefficient float64 `json:"coefficient"`
			URL         string  `json:"url"`
		} `json:"person_name_jaro_winkler"`
		PersonNamePhoneticMatch struct {
			Description string  `json:"description"`
			Coefficient float64 `json:"coefficient"`
			URL         string  `json:"url"`
		} `json:"person_name_phonetic_match"`
		NameFingerprintLevenshtein struct {
			Description string  `json:"description"`
			Coefficient float64 `json:"coefficient"`
			URL         string  `json:"url"`
		} `json:"name_fingerprint_levenshtein"`
		NameMetaphoneMatch struct {
			Description string `json:"description"`
			Coefficient int    `json:"coefficient"`
			URL         string `json:"url"`
		} `json:"name_metaphone_match"`
		NameSoundexMatch struct {
			Description string `json:"description"`
			Coefficient int    `json:"coefficient"`
			URL         string `json:"url"`
		} `json:"name_soundex_match"`
		AddressEntityMatch struct {
			Description string  `json:"description"`
			Coefficient float64 `json:"coefficient"`
			URL         string  `json:"url"`
		} `json:"address_entity_match"`
		CryptoWalletAddress struct {
			Description string  `json:"description"`
			Coefficient float64 `json:"coefficient"`
			URL         string  `json:"url"`
		} `json:"crypto_wallet_address"`
		IsinSecurityMatch struct {
			Description string  `json:"description"`
			Coefficient float64 `json:"coefficient"`
			URL         string  `json:"url"`
		} `json:"isin_security_match"`
		LeiCodeMatch struct {
			Description string  `json:"description"`
			Coefficient float64 `json:"coefficient"`
			URL         string  `json:"url"`
		} `json:"lei_code_match"`
		OgrnCodeMatch struct {
			Description string  `json:"description"`
			Coefficient float64 `json:"coefficient"`
			URL         string  `json:"url"`
		} `json:"ogrn_code_match"`
		VesselImoMmsiMatch struct {
			Description string  `json:"description"`
			Coefficient float64 `json:"coefficient"`
			URL         string  `json:"url"`
		} `json:"vessel_imo_mmsi_match"`
		InnCodeMatch struct {
			Description string  `json:"description"`
			Coefficient float64 `json:"coefficient"`
			URL         string  `json:"url"`
		} `json:"inn_code_match"`
		BicCodeMatch struct {
			Description string  `json:"description"`
			Coefficient float64 `json:"coefficient"`
			URL         string  `json:"url"`
		} `json:"bic_code_match"`
		IdentifierMatch struct {
			Description string  `json:"description"`
			Coefficient float64 `json:"coefficient"`
			URL         string  `json:"url"`
		} `json:"identifier_match"`
		WeakAliasMatch struct {
			Description string  `json:"description"`
			Coefficient float64 `json:"coefficient"`
			URL         string  `json:"url"`
		} `json:"weak_alias_match"`
		CountryMismatch struct {
			Description string  `json:"description"`
			Coefficient float64 `json:"coefficient"`
			URL         string  `json:"url"`
		} `json:"country_mismatch"`
		LastNameMismatch struct {
			Description string  `json:"description"`
			Coefficient float64 `json:"coefficient"`
			URL         string  `json:"url"`
		} `json:"last_name_mismatch"`
		DobYearDisjoint struct {
			Description string  `json:"description"`
			Coefficient float64 `json:"coefficient"`
			URL         string  `json:"url"`
		} `json:"dob_year_disjoint"`
		DobDayDisjoint struct {
			Description string  `json:"description"`
			Coefficient float64 `json:"coefficient"`
			URL         string  `json:"url"`
		} `json:"dob_day_disjoint"`
		GenderMismatch struct {
			Description string  `json:"description"`
			Coefficient float64 `json:"coefficient"`
			URL         string  `json:"url"`
		} `json:"gender_mismatch"`
		OrgidDisjoint struct {
			Description string  `json:"description"`
			Coefficient float64 `json:"coefficient"`
			URL         string  `json:"url"`
		} `json:"orgid_disjoint"`
		NumbersMismatch struct {
			Description string  `json:"description"`
			Coefficient float64 `json:"coefficient"`
			URL         string  `json:"url"`
		} `json:"numbers_mismatch"`
	} `json:"matcher"`
	Limit int `json:"limit"`
}

func Search() {
	// resp, err := http.Post(YENTE_URL, "application/json", )
}
