package api

import (
	"net/http"
	"encoding/json"
	"fmt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/blobstore"
)

type UploadUrl struct {
	Url string `json:"url"`
}

type PreviewUrl struct {
	Url string `json:"url"`
}

func init() {
	http.HandleFunc("/api/url", createUrl)
	http.HandleFunc("/api/callback", callback)
	http.HandleFunc("/api/video", serve)
}

func createUrl(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	option := blobstore.UploadURLOptions{
		MaxUploadBytes: 1024 * 1024 * 1024,
		StorageBucket:  "morita-demo",
	}

	url, err := blobstore.UploadURL(c, "/api/callback", &option)

	uploadUrl := UploadUrl{}
	uploadUrl.Url = url.String()

	outputJson, err := json.Marshal(&uploadUrl)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(outputJson))
}

func callback(w http.ResponseWriter, r *http.Request) {
	blobs, _, err := blobstore.ParseUpload(r)
	if err != nil {
		panic(err)
	}
	file := blobs["file"]
	if len(file) == 0 {
		panic(err)
	}

	previewUrl := PreviewUrl{}
	previewUrl.Url = "/api/video?key=" + string(file[0].BlobKey)

	outputJson, err := json.Marshal(&previewUrl)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(outputJson))
}

func serve(w http.ResponseWriter, r *http.Request) {
	blobKey := r.FormValue("key")
	blobstore.Send(w, appengine.BlobKey(blobKey))
}