### STAGE 1: Build ###
FROM golang:1.25.7-alpine3.22 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o /main

### STAGE 2: Run ###
FROM alpine:3.23.3

COPY --from=build /main ./main

EXPOSE 80

CMD [ "./main" ]
