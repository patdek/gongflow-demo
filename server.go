package main

import (
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/patdek/gongflow"
)

var (
	tempPath = path.Join(os.TempDir(), "gongflow")
)

func main() {
	// These assets are embedded into the executable using
	// https://github.com/jteeuwen/go-bindata
	// to generate bindata.go -- this removes the need to manage the path
	// just to run the demo ... it should "just work"(TM)
	listenAt := ":8080"
	bindStaticAssets()

	// ensure the tempPath exists
	os.MkdirAll(tempPath, 0777)

	// the actual demo, yey!
	go cleanupUploads() // loop forever doing cleanup
	http.HandleFunc("/upload", handleUpload)

	log.Println("Listening at:", listenAt)
	log.Fatal(http.ListenAndServe(listenAt, nil))
}

func cleanupUploads() {
	timeDur := time.Duration(1) * time.Minute
	t := time.NewTicker(timeDur)
	for _ = range t.C {
		err := gongflow.CleanupParts(tempPath, timeDur)
		if err != nil {
			log.Println(err)
		}
	}
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	fd, err := gongflow.ExtractFlowData(r) // extract the data that all flow uploads have
	if err != nil {
		http.Error(w, "Unable to extract flow data: "+err.Error(), 500)
	}

	if r.Method == "GET" { // status request
		msg, code := gongflow.CheckPart(tempPath, fd)
		http.Error(w, msg, code)
	} else if r.Method == "POST" { // upload part
		msg, err := gongflow.UploadPart(tempPath, fd, r)
		if err != nil {
			http.Error(w, "WAT: "+err.Error(), 500)
			return
		}
		if msg != "" {
			http.Error(w, msg+" is done", 200)
			return
		}
		http.Error(w, "continuing", 200)
	} else { // bug
		http.Error(w, "Bad Method", 500)
		return
	}

}

func bindStaticAssets() {
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
