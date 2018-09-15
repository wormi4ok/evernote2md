FROM golang:1.11 AS build
# Modules are enabled to use exact versions of dependencies
ENV GO111MODULE on

WORKDIR /go/src/github.com/wormi4ok/evernote2md

COPY . .

RUN go install && go test ./...

FROM alpine:3.8

LABEL   org.label-schema.name="evernote2md" \
        org.label-schema.description="Convert Evernote .enex export file to Markdown" \
        org.label-schema.vcs-url="https://github.com/wormi4ok/evernote2md" \
        org.label-schema.docker.cmd="docker run --rm wormi4ok/evernote2md export.enex notes"

COPY --from=build /go/bin/evernote2md /

ENTRYPOINT ["/evernote2md"]

CMD [ "-h" ]