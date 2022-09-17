############################
# STEP 1 build executable binary
############################
FROM golang:1.19-alpine as builder

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o exec -ldflags="-s -w" ./...

############################
# STEP 2 build a small image
############################
FROM scratch

COPY --from=builder /app/exec /app/exec

ENTRYPOINT [ "/app/exec" ]
