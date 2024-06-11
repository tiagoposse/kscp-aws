FROM golang as build

ENV CGO_ENABLED=0

ADD main.go go.mod go.sum
RUN go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o /agent *.go


FROM scratch

COPY --from=build /agent /agent

ENTRYPOINT [ "/agent" ]