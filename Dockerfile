FROM golang as build

ENV CGO_ENABLED=0

ADD main.go go.mod go.sum .
RUN go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o /agent *.go

FROM golang as certs

RUN echo | openssl s_client -showcerts -servername sts.amazonaws.com -connect sts.amazonaws.com:443 2>/dev/null | openssl x509 -outform PEM > /usr/local/share/ca-certificates/sts.amazonaws.com.crt
RUN update-ca-certificates

FROM scratch

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /agent /agent

ENTRYPOINT [ "/agent" ]