package yenti

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"bitbucket.org/msafaridanquah/verifylab-service/foundation/envvar"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/logger"
)

type Config struct {
	Env *envvar.Configuration
	Log *logger.Logger
}

type Yenti struct {
	endpoint string
	http     *http.Client
	log      *logger.Logger
}
type Weights struct {
	NameLiteralMatch float64 `json:"name_literal_match"`
	NameSoundexMatch float64 `json:"name_soundex_match"`
}

type Query struct {
	Properties Properties `json:"properties"`
	Schema     string     `json:"schema"`
}

type Properties struct {
	Name               []string `json:"name"`
	Nationality        []string `json:"nationality,omitempty"`
	BirthDate          []string `json:"birthDate,omitempty"`
	Jurisdiction       []string `json:"jurisdiction,omitempty"`
	RegistrationNumber []string `json:"registrationNumber,omitempty"`
}

type ResponsePropertyMatcher struct {
	Description string  `json:"description"`
	Coefficient float32 `json:"coefficient"`
	URL         string  `json:"url"`
}

type NewLookup struct {
	Weights Weights          `json:"weights"`
	Queries map[string]Query `json:"queries"`
}

type ResponseProperty struct {
	Status  int `json:"status"`
	Results []struct {
		ID         string `json:"id"`
		Caption    string `json:"caption"`
		Schema     string `json:"schema"`
		Properties struct {
			Name []string `json:"name"`
		} `json:"properties"`
		Datasets   []string       `json:"datasets"`
		Referents  []string       `json:"referents"`
		Target     bool           `json:"target"`
		FirstSeen  string         `json:"first_seen"`
		LastSeen   string         `json:"last_seen"`
		LastChange string         `json:"last_change"`
		Score      float64        `json:"score"`
		Features   map[string]any `json:"features"`
		Match      bool           `json:"match"`
		Token      string         `json:"token"`
	} `json:"results"`
	Total struct {
		Value    int    `json:"value"`
		Relation string `json:"relation"`
	} `json:"total"`
	Query struct {
		ID         string `json:"id"`
		Schema     string `json:"schema"`
		Properties struct {
			Name []string `json:"name"`
		} `json:"properties"`
	} `json:"query"`
}

type YentiResponse struct {
	Properties map[string]ResponseProperty        `json:"responses"`
	Matcher    map[string]ResponsePropertyMatcher `json:"matcher"`
	Limit      int                                `json:"limit"`
}

var defaultClient = http.Client{
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 15 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
}

func New(conf Config) (*Yenti, error) {
	get := func(v string) string {
		res, err := conf.Env.Get(v)
		if err != nil {
			conf.Log.Error(context.Background(), "env failed")
		}

		return res
	}

	host := get("YENTE_HOST")

	return &Yenti{
		endpoint: strings.Join([]string{host, "/match/default?algorithm=best"}, ""),
		http:     &defaultClient,
		log:      conf.Log,
	}, nil
}

func (o *Yenti) Search(nu NewLookup) (YentiResponse, error) {
	payload, err := json.Marshal(nu)
	if err != nil {
		return YentiResponse{}, err
	}

	u, err := url.Parse(o.endpoint)
	if err != nil {
		return YentiResponse{}, fmt.Errorf("parsing endpoint: %w", err)
	}

	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(payload))
	if err != nil {
		o.log.Error(context.Background(), "Error creating request: %v", err)
	}

	// Create an HTTP client and send the request
	resp, err := o.http.Do(req)
	if err != nil {
		o.log.Error(context.Background(), "Error sending request: %v", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			return
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		o.log.Error(context.Background(), "Error parsing response: %v", err)
	}

	var yentiResponse YentiResponse

	if err := json.Unmarshal(body, &yentiResponse); err != nil {
		return YentiResponse{}, err
	}

	return yentiResponse, nil
}
