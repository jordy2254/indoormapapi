FROM golang:latest
LABEL maintainer="jhellier <jordanhellier@googlemail.com>"
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
ENV PORT 8081
RUN go build cmd/restservice/restservice.go
CMD ["./restservice"]