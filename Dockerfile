
#Build stage
FROM golang:1.19.0-alpine as build-env
 
ENV APP_NAME fn-kube-state
ENV CMD_PATH /cmd/main.go
 
COPY . $GOPATH/src/$APP_NAME
WORKDIR $GOPATH/src/$APP_NAME
 
RUN CGO_ENABLED=0 go build -v -o /$APP_NAME $GOPATH/src/$APP_NAME/$CMD_PATH
 
# Run Stage
FROM alpine:3.14
 
ENV APP_NAME fn-kube-state
ENV KUBE_CLIENT inCluster
 
COPY --from=build-env /$APP_NAME .
 
EXPOSE 8080

# Start app
CMD ./$APP_NAME