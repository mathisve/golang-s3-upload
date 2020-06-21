package main

import (
  "fmt"
  "html/template"
  "net/http"
  "math/rand"
  "encoding/json"
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
   sessions = make(map[string]*session.Session)
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

  if err != nil {
    fmt.Println(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  region := getBucketRegion(bucket)
  uploader := s3manager.NewUploader(awsConnectRegion(region))
  _, err = uploader.Upload(&s3manager.UploadInput{
    Bucket: aws.String(bucket),
    Key:    aws.String(strconv.Itoa(rand.Int())[:5]+handler.Filename),
    Body:   file,
    ACL: aws.String("public-read"),
  })
  if err != nil {
    fmt.Println(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
}

func getBucketObjects(w http.ResponseWriter, r* http.Request) {
    r.ParseMultipartForm(128 << 20)

    bucket := r.Form["bucket"][0]
    bo := listBucketItems(bucket)
    json.NewEncoder(w).Encode(bo)
}

func awsConnectRegion(region string) *session.Session {
  if region == "" {
    region = aws_region
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
      fmt.Println(err)
    }
    sessions[region] = sess
    return sess
  }
}

func gets3clientRegion(region string) *s3.S3 {
  if region == "" {
    region = aws_region
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
    fmt.Println(err)
  }

  for _, bucket := range(result.Buckets) {
    b = append(b, Bucket {
      Name: *bucket.Name,
    })
  }

  return b
}

func getBucketRegion(bucket string) (region string){
  url := fmt.Sprintf("https://%s.s3.amazonaws.com", bucket)
  res, err := http.Head(url)
  if err != nil {
     fmt.Println(err)
  }
  return res.Header.Get("X-Amz-Bucket-Region")
}

func listBucketItems(bucket string) (bo []BucketObject) {
  region := getBucketRegion(bucket)
  resp, err := gets3clientRegion(region).ListObjectsV2(&s3.ListObjectsV2Input {
    Bucket: aws.String(bucket),
  })
  if err != nil {
    fmt.Println(err)
  }

  for _, object := range(resp.Contents) {
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

  err := http.ListenAndServe(":" + port, router)
  if err != nil {
    panic(err)
  }
}

func main() {
  fmt.Println("localhost:" + port)

  router()
}
