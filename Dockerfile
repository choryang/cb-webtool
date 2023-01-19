FROM golang:1.19.4 as prod

WORKDIR /cb-webtool 
COPY . .

#RUN apk update && apk add git
#RUN apk add --no-cache bash
RUN apk update
RUN apk add --no-cache bash git gcc

#RUN go get -u -v github.com/go-session/echo-session
#RUN go get -u github.com/labstack/echo/...
#RUN go get -u github.com/davecgh/go-spew/spew
ENV GO111MODULE on
RUN go get github.com/cespare/reflex@latest

EXPOSE 1235

CMD reflex -r '\.(html|go)' -s go run main.go
