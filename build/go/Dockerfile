FROM golang

RUN go get github.com/dgrijalva/jwt-go
RUN go get github.com/gorilla/mux
RUN go get gopkg.in/redis.v5
RUN go get github.com/parnurzeal/gorequest
RUN go get github.com/gorilla/websocket

# Expose the application on port 80
EXPOSE 80

RUN mkdir -p /go/src/app
WORKDIR /go/src/app

# this will ideally be built by the ONBUILD below ;)
CMD ["run"]