package main

import (
  "fmt"
  "html/template"
  "net/http"
  "math/rand"
  "strconv"
  "github.com/gorilla/mux"

  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/s3"
  "github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const (
	staticDir = "/static"
  port = "8080"

  aws_region = "eu-central-1"
)

var (
   sess = awsConnectRegion("")
   svc = gets3client()
   sessions = make(map[string]*session.Session)
)

type IndexPageData struct {
  Buckets []Bucket
}

type Bucket struct {
  Name string
}

func index(w http.ResponseWriter, r* http.Request) {
  tmpl := template.Must(template.ParseFiles("static/index.html"))

  data := IndexPageData{
    Buckets: getBuckets(),
  }

  tmpl.Execute(w, data)
}

func upload(w http.ResponseWriter, r* http.Request) {
  r.ParseMultipartForm(128 << 20)

  bucket := r.Form["bucket"][0]

  file, handler, err := r.FormFile("file")
  if err != nil {
    fmt.Println(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  defer file.Close()

  fmt.Printf("UploadingFile: %+v to %+v\n", handler.Filename, bucket)

  //get region of bucket
  url := fmt.Sprintf("https://%s.s3.amazonaws.com", bucket)
  res, err := http.Head(url)
  if err != nil {
     panic(err)
  }
  region := res.Header.Get("X-Amz-Bucket-Region")

  if err != nil {
    panic(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  uploader := s3manager.NewUploader(awsConnectRegion(region))
  _, err = uploader.Upload(&s3manager.UploadInput{
    Bucket: aws.String(bucket),
    Key:    aws.String(strconv.Itoa(rand.Int())[:5]+handler.Filename),
    Body:   file,
    ACL: aws.String("public-read"),
  })
  if err != nil {
    panic(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
}

func awsConnectRegion(region string) *session.Session {
  if region == "" {
    region = "eu-central-1"
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
      panic(err)
    }
    sessions[region] = sess
    return sess
  }
}

func gets3client() *s3.S3 {
  return s3.New(sess)
}

func getBuckets() (b []Bucket) {
  result, err := svc.ListBuckets(&s3.ListBucketsInput{})
  if err != nil {
    panic(err)
  }

  for _, bucket := range(result.Buckets) {
    b = append(b, Bucket {Name: *bucket.Name})
  }

  return b
}

func router() {
  router := mux.NewRouter()
  router.PathPrefix(staticDir).Handler(http.StripPrefix(staticDir, http.FileServer(http.Dir("."+staticDir))))

  router.HandleFunc("/", index).Methods("GET")
  router.HandleFunc("/upload", upload).Methods("POST")
  err := http.ListenAndServe(":"+port, router)
  if err != nil {
    panic(err)
  }
}

func main() {
  fmt.Println("Server online!")
  router()
}
