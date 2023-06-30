FROM golang:alpine as app-builder
WORKDIR $GOPATH/src/ticket-pimp/
COPY . .
RUN apk add git
# Static build required so that we can safely copy the binary over.
# `-tags timetzdata` embeds zone info from the "time/tzdata" package.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -ldflags '-extldflags "-static"' -tags timetzdata cmd/main.go

FROM scratch
# the test program:
COPY --from=app-builder /go/bin/main /ticket-pimp
COPY --from=app-builder /go/src/ticket-pimp/cmd/.env /
# the tls certificates:
# NB: this pulls directly from the upstream image, which already has ca-certificates:
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/ticket-pimp"]