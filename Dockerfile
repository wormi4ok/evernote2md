FROM golang:1.20-alpine3.18 AS build

ENV CGO_ENABLED 0

WORKDIR /go/src/github.com/wormi4ok/evernote2md

RUN set -xe && apk add --no-cache git

COPY . .

RUN go test ./... && go install -trimpath -ldflags "-s -w -X main.version=$(git describe --tags --abbrev=0)"

FROM alpine:3.18

LABEL   org.opencontainers.image.title="evernote2md" \
        org.opencontainers.image.description="Convert Evernote .enex export file to Markdown" \
        org.opencontainers.image.source="https://github.com/wormi4ok/evernote2md" \
        org.opencontainers.image.authors="wormi4ok"

COPY --from=build /go/bin/evernote2md /

ENTRYPOINT ["/evernote2md"]

CMD [ "-h" ]
