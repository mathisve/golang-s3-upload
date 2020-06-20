package main

import (
  "fmt"
  "net/http"
  "math/rand"
  "strconv"
  "github.com/gorilla/mux"

  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  //"github.com/aws/aws-sdk-go/service/s3"
  "github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const (
	staticDir = "/static"
  port = "8080"

  aws_region = "eu-central-1"
  aws_bucket = "golang-s3-uploader"
)

var sess = awsconnect()

func index(w http.ResponseWriter, r* http.Request) {
  http.ServeFile(w, r, "static/index.html")
}

func upload(w http.ResponseWriter, r* http.Request) {
  r.ParseMultipartForm(30 << 20)

  file, handler, err := r.FormFile("file")
  if err != nil {
    fmt.Println(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
  defer file.Close()

  fmt.Printf("UploadingFile: %+v\n", handler.Filename)

  uploader := s3manager.NewUploader(sess)
  _, err = uploader.Upload(&s3manager.UploadInput{
    Bucket: aws.String(aws_bucket),
    Key:    aws.String(strconv.Itoa(rand.Int())[:5]+handler.Filename),
    Body:   file,
    ACL: aws.String("public-read"),
  })
  if err != nil {
    panic(err)
    return
  }
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

func awsconnect() *session.Session {
  sess, err := session.NewSession(
    &aws.Config{
      Region: aws.String(aws_region),
    },
  )
  if err != nil {
    panic(err)
  }
  return sess
}

func main() {
  fmt.Println("Server online!")
  router()
}
