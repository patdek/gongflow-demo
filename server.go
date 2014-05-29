package main

import (
	"log"
	"net/http"

	"github.com/patdek/ngflow"
)

func main() {
	// These assets are embedded into the executable using
	// https://github.com/jteeuwen/go-bindata
	// to generate bindata.go -- this removes the need to manage the path
	// just to run the demo ... it should "just work"(TM)
	angular, ng, bootstrap, app, index, glyphicons := mustLoadAssets()
	http.HandleFunc("/angular.min.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		w.Write(angular)
	})
	http.HandleFunc("/ng-flow-standalone.min.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		w.Write(ng)
	})
	http.HandleFunc("/app.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		w.Write(app)
	})
	http.HandleFunc("/bootstrap-combined.min.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		w.Write(bootstrap)
	})
	http.HandleFunc("/img/glyphicons-halflings.png", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.Write(glyphicons)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(index)
	})
	http.HandleFunc("/upload", ngflow.UploadHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func mustLoadAssets() ([]byte, []byte, []byte, []byte, []byte, []byte) {
	angular, err := Asset("html/angular.min.js")
	if err != nil {
		log.Fatal("Problem loading angular static asset: ", err)
	}
	ng, err := Asset("html/ng-flow-standalone.min.js")
	if err != nil {
		log.Fatal("Problem loading ng-flow static asset: ", err)
	}
	bootstrap, err := Asset("html/bootstrap-combined.min.css")
	if err != nil {
		log.Fatal("Problem loading bootstrap static asset: ", err)
	}
	app, err := Asset("html/app.js")
	if err != nil {
		log.Fatal("Problem loading app static asset: ", err)
	}
	index, err := Asset("html/index.html")
	if err != nil {
		log.Fatal("Problem loading index.html static asset: ", err)
	}
	glyphicons, err := Asset("html/glyphicons-halflings.png")
	if err != nil {
		log.Fatal("Problem loading glyphicons static asset: ", err)
	}
	return angular, ng, bootstrap, app, index, glyphicons
}
