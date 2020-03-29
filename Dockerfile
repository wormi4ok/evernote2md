FROM golang:1.14-alpine AS build

ENV CGO_ENABLED 0

WORKDIR /go/src/github.com/wormi4ok/evernote2md

COPY . .

RUN set -xe && apk add --no-cache git

RUN go install && go test ./...

FROM alpine:3.11

LABEL   org.label-schema.name="evernote2md" \
        org.label-schema.description="Convert Evernote .enex export file to Markdown" \
        org.label-schema.vcs-url="https://github.com/wormi4ok/evernote2md" \
        org.label-schema.docker.cmd="docker run --rm wormi4ok/evernote2md export.enex notes"

COPY --from=build /go/bin/evernote2md /

ENTRYPOINT ["/evernote2md"]

CMD [ "-h" ]
