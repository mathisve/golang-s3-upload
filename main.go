package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	url2 "net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const (
	staticDir = "/static"
	localhost = "127.0.0.1"
	port      = "80"

	awsRegion = "eu-central-1"
)

var (
	sessions   = make(map[string]*session.Session)
	s3sessions = make(map[string]*s3.S3)
)

type IndexPageData struct {
	Buckets []Bucket
}

type Bucket struct {
	Name string
}

type BucketObject struct {
	Name string
	Size int64
}

func index(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("static/index.html"))

	data := IndexPageData{
		Buckets: getBuckets(),
	}

	err := tmpl.Execute(w, data)
	if err != nil {
		log.Println(err)
	}
}

func upload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(128 << 20)
	if err != nil {
		log.Println(err)
	}

	bucket := r.Form["bucket"][0]

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer file.Close()

	log.Printf("UploadingFile: %+v to %+v\n", handler.Filename, bucket)

	region := getBucketRegion(bucket)
	uploader := s3manager.NewUploader(awsConnectRegion(region))
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(handler.Filename),
		Body:   file,
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			log.Println(awsErr)
		} else {
			log.Println(err)
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("Upload of file %v done!\n", handler.Filename)
}

func getBucketObjects(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(128 << 20)
	if err != nil {
		log.Println(err)
	}

	bucket := r.Form["bucket"][0]
	bo := listBucketItems(bucket)
	json.NewEncoder(w).Encode(bo)
}

func awsConnectRegion(region string) *session.Session {
	if region == "" {
		region = awsRegion
	}

	if val, ok := sessions[region]; ok {
		return val
	} else {
		sess, err := session.NewSession(
			&aws.Config{
				Region: aws.String(region),
			},
		)
		if err != nil {
			log.Println(err)
		}
		sessions[region] = sess
		return sess
	}
}

func gets3clientRegion(region string) *s3.S3 {
	if region == "" {
		region = awsRegion
	}

	if val, ok := s3sessions[region]; ok {
		return val
	} else {
		s := s3.New(awsConnectRegion(region))
		s3sessions[region] = s
		return s
	}
}

func getBuckets() (b []Bucket) {
	result, err := gets3clientRegion("").ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		log.Println(err)
	}

	for _, bucket := range result.Buckets {
		b = append(b, Bucket{
			Name: *bucket.Name,
		})
	}

	return b
}

func getBucketRegion(bucket string) (region string) {
	url := fmt.Sprintf("https://%s.s3.amazonaws.com", url2.QueryEscape(bucket))
	res, err := http.Head(url)
	if err != nil {
		log.Println(err)
	}
	return res.Header.Get("X-Amz-Bucket-Region")
}

func listBucketItems(bucket string) (bo []BucketObject) {
	region := getBucketRegion(bucket)
	resp, err := gets3clientRegion(region).ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		log.Println(err)
	}

	for _, object := range resp.Contents {
		bo = append(bo, BucketObject{
			Name: *object.Key,
			Size: *object.Size,
		})
	}

	return bo
}

func router() {
	router := mux.NewRouter()
	router.PathPrefix(staticDir).Handler(http.StripPrefix(staticDir, http.FileServer(http.Dir("."+staticDir))))

	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/upload", upload).Methods("POST")
	router.HandleFunc("/getBucketObjects", getBucketObjects).Methods("POST")

	log.Printf("Now live on http://%v:%v", localhost, port)

	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		panic(err)
	}
}

func main() {
	router()
}
