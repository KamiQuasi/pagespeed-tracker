package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/quad"
	// "github.com/google/go-github/github"
)

// Site struct
type Site struct {
	URL        string
	Protocol   string
	Repository string
	Analysis   []Analyzer
	Scores     []PageSpeedResult
}

func (p *Site) analyze() error {
	request := PageSpeedRequest{}
	request.Strategy = "mobile"
	p.Scores = append(p.Scores, request.get(p)[0])

	return nil
}

// Analyzer interface
type Analyzer interface {
	analyze()
}

// func loadSite(site string) (*Site, error) {
// 	store, err := cayley.NewMemoryGraph()
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	p := cayley.StartPath(store, quad.String(site))

// 	err = p.Iterate(nil).EachValue(nil, func(value quad.Value) {
// 		nativeValue := quad.NativeOf(value)

// 	})

// 	store.AddQuad(quad.Make("phrase of the day", "is of course", "Hello World!", nil))

// 	// Now we iterate over results. Arguments:
// 	// 1. Optional context used for cancellation.
// 	// 2. Flag to optimize query before execution.
// 	// 3. Quad store, but we can omit it because we have already built path with it.
// 	err = p.Iterate(nil).EachValue(nil, func(value quad.Value) {
// 		nativeValue := quad.NativeOf(value) // this converts RDF values to normal Go types
// 		fmt.Println(nativeValue)
// 	})
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// }

// Repository struct
type Repository struct {
	URL string
}

// PageSpeedRequest struct
type PageSpeedRequest struct {
	URL                       string
	FilterThirdPartyResources bool
	Locale                    string
	Rule                      string
	Screenshot                bool
	Strategy                  string
}

func (ps *PageSpeedRequest) get(sites ...*Site) []PageSpeedResult {
	var results []PageSpeedResult

	for _, site := range sites {
		resp, err := http.Get(fmt.Sprintf("https://www.googleapis.com/pagespeedonline/v2/runPagespeed?url=%s&strategy=%s", site.URL, ps.Strategy))
		if err != nil {

		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {

		}
		result := PageSpeedResult{}
		json.Unmarshal(body, &result)
		results = append(results, result)
	}

	return results
}

// PageSpeedResult struct
type PageSpeedResult struct {
	Kind             string              `json:"kind"`
	ID               string              `json:"id"`
	ResponseCode     int                 `json:"responseCode"`
	Title            string              `json:"title"`
	RuleGroups       PageSpeedRuleGroups `json:"ruleGroups,omitempty"`
	PageStats        PageStats           `json:"pageStats"`
	Version          Version             `json:"version"`
	FormattedResults FormattedResults    `json:"formattedResults"`
}

// PageSpeedRuleGroups struct
type PageSpeedRuleGroups struct {
	Speed     PageSpeedRule `json:"SPEED,omitempty"`
	Usability PageSpeedRule `json:"USABILITY,omitempty"`
}

// PageSpeedRule struct
type PageSpeedRule struct {
	Score int `json:"score,omitempty"`
}

// PageStats struct
type PageStats struct {
	NumberResources         int `json:"numberResources"`
	NumberHosts             int `json:"numberHosts"`
	TotalRequestBytes       int `json:"totalRequestBytes,string"`
	NumberStaticResources   int `json:"numberStaticResources,string"`
	HTMLResponseBytes       int `json:"htmlResponseBytes,string"`
	CSSResponseBytes        int `json:"cssResponseBytes,string"`
	ImageResponseBytes      int `json:"imageResponseBytes,string"`
	JavaScriptResponseBytes int `json:"javascriptResponseBytes,string"`
	OtherResponseBytes      int `json:"otherResponseBytes,string"`
	NumberJsResources       int `json:"numberJsResources"`
	NumberCSSResorces       int `json:"numberCssResources"`
}

// Version struct
type Version struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
}

// FormattedResults struct
type FormattedResults struct {
	Locale      string      `json:"locale"`
	RuleResults RuleResults `json:"ruleResults"`
}

// RuleResults struct
type RuleResults struct {
	AvoidLandingPageRedirects       Rule `json:"AvoidLandingPageRedirects"`
	EnableGzipCompression           Rule `json:"EnableGzipCompression"`
	LeverageBrowserCaching          Rule `json:"LeverageBrowserCaching"`
	MainResourceServerResponseTime  Rule `json:"MainResourceServerResponseTime"`
	MinifyCSS                       Rule `json:"MinifyCss"`
	MinifyHTML                      Rule `json:"MinifyHTML"`
	MinifyJavaScript                Rule `json:"MinifyJavaScript"`
	MinimizeRenderBlockingResources Rule `json:"MinimizeRenderBlockingResources"`
	OptimizeImages                  Rule `json:"OptimizeImages"`
	PrioritizeVisibleContent        Rule `json:"PrioritizeVisibleContent"`
}

// Rule struct
type Rule struct {
	LocalizedRuleName string     `json:"localizedRuleName"`
	RuleImpact        float64    `json:"ruleImpact"`
	Groups            []string   `json:"groups"`
	Summary           Formatter  `json:"summary"`
	URLBlocks         []URLBlock `json:"urlBlocks"`
}

// URLBlock struct
type URLBlock struct {
	Header Formatter   `json:"header"`
	URLs   []URLResult `json:"urls"`
}

// URLResult struct
type URLResult struct {
	Result Formatter `json:"result"`
}

// Formatter struct
type Formatter struct {
	Format string     `json:"format"`
	Args   []Argument `json:"args"`
}

// Argument struct
type Argument struct {
	Type  string `json:"type"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

func main() {
	http.HandleFunc("/site/", siteHandler)
	http.HandleFunc("/", home)
	bind := fmt.Sprintf("%s:%s", os.Getenv("OPENSHIFT_GO_IP"), os.Getenv("OPENSHIFT_GO_PORT"))
	fmt.Printf("listening on %s...", bind)
	err := http.ListenAndServe(bind, nil)
	if err != nil {
		panic(err)
	}
}

func siteHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path[len("/site/"):]
	site := &Site{URL: url}
	t, _ := template.ParseFiles("templates/site.html")
	t.Execute(w, site)
}
func home(res http.ResponseWriter, req *http.Request) {
	siteURL := "http://developers.redhat.com"
	// Create a brand new graph
	store, err := cayley.NewMemoryGraph()
	if err != nil {
		log.Fatalln(err)
	}

	store.AddQuad(quad.Make(siteURL, "type", "site", nil))
	store.AddQuad(quad.Make(siteURL, "name", "Red Hat Developers", nil))
	store.AddQuad(quad.Make(siteURL, "allows protocol", "http", nil))
	store.AddQuad(quad.Make(siteURL, "allows protocol", "https", nil))
	store.AddQuad(quad.Make(siteURL, "scores", 72, nil))
	// Now we create the path, to get to our data
	p := cayley.StartPath(store, quad.String(siteURL)).Out()

	// Now we iterate over results. Arguments:
	// 1. Optional context used for cancellation.
	// 2. Flag to optimize query before execution.
	// 3. Quad store, but we can omit it because we have already built path with it.
	err = p.Iterate(nil).EachValue(nil, func(value quad.Value) {
		nativeValue := quad.NativeOf(value) // this converts RDF values to normal Go types
		fmt.Println(nativeValue)
	})
	if err != nil {
		log.Fatalln(err)
	}

	// client := github.NewClient(nil)

	// orgs, _, err := client.Organizations.List("KamiQuasi", nil)

	developers := Site{siteURL, "http", "https://github.com/redhat-developer/developers.redhat.com", nil, nil}
	developers.analyze()

	b, err := json.Marshal(developers.Scores[0])
	if err != nil {

	}
	fmt.Fprintf(res, "%s", b)
}
