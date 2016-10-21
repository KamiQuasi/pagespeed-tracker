package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// PageSpeed struct
type PageSpeed struct {
	Kind         string    `json:"numberResources"`
	ID           string    `json:"numberResources"`
	ResponseCode int       `json:"numberResources"`
	Title        string    `json:"numberResources"`
	RuleGroups   RuleGroup `json:"ruleGroups"`
	PageStats    PageStats `json:"pageStats"`
	Version      Version   `json:"Version"`
}

// RuleGroup struct
type RuleGroup struct {
	Speed     Rule `json:"SPEED"`
	Usability Rule `json:"USABILITY"`
}

// Rule struct
type Rule struct {
	Score int `json:"score"`
}

// PageStats struct
type PageStats struct {
	NumberResources         int    `json:"numberResources"`
	NumberHosts             int    `json:"numberHosts"`
	TotalRequestBytes       string `json:"totalRequestBytes"`
	NumberStaticResources   int    `json:"numberStaticResources"`
	HTMLResponseBytes       string `json:"htmlResponseBytes"`
	CSSResponseBytes        string `json:"cssResponseBytes"`
	ImageResponseBytes      string `json:"imageResponseBytes"`
	JavaScriptResponseBytes string `json:"javascriptResponseBytes"`
	OtherResponseBytes      string `json:"otherResponseBytes"`
	NumberJsResources       int    `json:"numberJsResources"`
	NumberCSSResorces       int    `json:"numberCssResources"`
}

// Version struct
type Version struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
}

func main() {
	http.HandleFunc("/", home)
	bind := fmt.Sprintf("%s:%s", os.Getenv("OPENSHIFT_GO_IP"), os.Getenv("OPENSHIFT_GO_PORT"))
	fmt.Printf("listening on %s...", bind)
	err := http.ListenAndServe(bind, nil)
	if err != nil {
		panic(err)
	}
}

func home(res http.ResponseWriter, req *http.Request) {
	resp, err := http.Get("https://www.googleapis.com/pagespeedonline/v2/runPagespeed?url=http://developers.redhat.com")
	if err != nil {

	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {

	}
	result := PageSpeed{}
	json.Unmarshal(body, &result)
	enc := json.NewEncoder(os.Stdout)
	fmt.Fprintf(res, "%s", enc.Encode(result))
}
