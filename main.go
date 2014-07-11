package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	cfenv "github.com/cloudfoundry-community/go-cfenv"
	"github.com/go-martini/martini"
	"github.com/longnguyen11288/elastigo/api"
	"github.com/longnguyen11288/elastigo/core"
	"github.com/longnguyen11288/elastigo/indices"
)

func AddData(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	core.Index("gophers", "data", "", nil, string(body))
	fmt.Fprint(w, `{ "success": "true" }`)
}

func main() {
	appEnv, enverr := cfenv.Current()
	if enverr != nil {
		api.Host = "http://localhost:9200"
	} else {
		elasticSearch, err := appEnv.Services.WithTag("elasticsearch")
		if err == nil {
			api.Host = elasticSearch[0].Credentials["uri"]
		} else {
			log.Fatal("Unable to find elastic search service")
		}
	}
	indices.Create("gophers")
	m := martini.Classic()
	m.Post("/", AddData)
	m.Run()
}
