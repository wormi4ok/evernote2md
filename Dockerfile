FROM golang:1.16-alpine3.13 AS build

ENV CGO_ENABLED 0

WORKDIR /go/src/github.com/wormi4ok/evernote2md

RUN set -xe && apk add --no-cache git

COPY . .

RUN go test ./... && go install -trimpath -ldflags "-s -w -X main.version=$(git describe --tags --abbrev=0)"

FROM alpine:3.13

LABEL   org.label-schema.name="evernote2md" \
        org.label-schema.description="Convert Evernote .enex export file to Markdown" \
        org.label-schema.vcs-url="https://github.com/wormi4ok/evernote2md" \
        org.label-schema.docker.cmd="docker run --rm wormi4ok/evernote2md export.enex notes"

COPY --from=build /go/bin/evernote2md /

ENTRYPOINT ["/evernote2md"]

CMD [ "-h" ]
