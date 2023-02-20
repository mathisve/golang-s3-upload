[![Dockerhub](https://github.com/mathisve/golang-s3-upload/actions/workflows/dockerpush.yaml/badge.svg)](https://github.com/mathisve/golang-s3-upload/actions/workflows/dockerpush.yaml)

# golang-s3-upload
Simple web tool to upload files to s3 buckets in bulk.

## How to use:
### Build yourself

Firstly build the docker container

`docker build -t s3upload .`

Then run it

`docker run -it -p 80:80 -e AWS_ACCESS_KEY_ID={YOUR ID} -e AWS_SECRET_ACCESS_KEY={YOUR KEY} s3upload`

### Use pre-built  container

`docker run -it -p 80:80 -e AWS_ACCESS_KEY_ID={YOUR ID} -e AWS_SECRET_ACCESS_KEY={YOUR KEY} mathisve/s3upload`
