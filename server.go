// Package gongflow-demo provides an instant-on working demo of using the gongflow library.  See
// http://godoc.org/github.com/patdek/gongflow for more information.
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

func init() {
	// ensure the tempPath exists
	os.MkdirAll(tempPath, 0777)
}

func main() {
	listenAt := ":8080" // listen on all local addresses :8080
	mapStaticAssets()

	// the actual demo, yey!
	go cleanupUploads() // loop forever doing cleanup
	http.HandleFunc("/upload", uploadHandler)

	log.Println("Listening at:", listenAt)
	log.Fatal(http.ListenAndServe(listenAt, nil))
}

// cleanupUploads is an example of how to write a loop to cleanup your temporary directory,
// it is a lazy implementation.  You should pass in a channel to signal it to close.
func cleanupUploads() {
	loopDur := time.Duration(1) * time.Minute   // loop every minute
	tooOldDur := time.Duration(5) * time.Minute // older than 5 minutes to be deleted
	t := time.NewTicker(loopDur)
	// this will "tick" every loopDur forever.
	for _ = range t.C {
		err := gongflow.PartsCleanup(tempPath, tooOldDur) // delete stuff in tempPath older than tooOldDur
		if err != nil {
			log.Println(err)
		}
	}
}

// uploadHandler is an example of how to write a handler for the two type of requests ng-flow
// will send.  It sends POST and GET requests.  POST to do the actual upload, GET to ask for
// status on parts.  See the ng-flow docs for more information.
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	ngFlowData, err := gongflow.PartFlowData(r)
	if err != nil {
		http.Error(w, "Unable to extract ngFlowData: "+err.Error(), 500)
	}

	if r.Method == "GET" { // ng-flow status request
		msg, code := gongflow.PartStatus(tempPath, ngFlowData)
		http.Error(w, msg, code)
	} else if r.Method == "POST" { // ng-flow upload part
		filePath, err := gongflow.PartUpload(tempPath, ngFlowData, r)
		if err != nil {
			http.Error(w, "Part Upload Failure: "+err.Error(), 500)
			return
		}
		if filePath != "" {
			http.Error(w, filePath+" is done", 200)

			// TODO: Add what you want to do with the file here, you want to get
			// it out of the temporary directory before any cleanup code you might
			// have written is run (like the clenaupUploads above
			log.Println("Part Upload Done: " + filePath)

			return
		}
		http.Error(w, "continuing to upload parts", 200)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// below here is just busywork about serving static assets
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// mapStaticAssets maps the various endpoints needed for the demo excluding the
// actual upload handler
func mapStaticAssets() {
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

// mustLoadAssets uses the bindata (see README.md) to load the assets so
// it doesn't have to worry about puzzling out the path dynamically.
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
