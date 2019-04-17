FROM golang:1.12-alpine
## Installs rsync and git
RUN apk add --no-cache \
    rsync \
    git
## Installs dependencies (listed in get.sh file)
WORKDIR /go/src/app
COPY get.sh ./
RUN sh get.sh
## Installs codebase and makes sure all dependencies are there.
COPY . .
RUN go get ./... \
  && go install ./...
EXPOSE 2112