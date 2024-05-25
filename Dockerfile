FROM golang:1.22.0-alpine3.19 as dependencies
COPY go.mod go.sum ./
RUN go mod download

FROM dependencies AS build
COPY . ./
RUN apk add build-base
RUN CGO_ENABLED=1 go build -o /main -ldflags="-w -s" .

FROM tihmmm/golang-alpine-rootless:go-1.22.0-alp-3.19
COPY --from=build /main /home/user/main
CMD ["./main"]