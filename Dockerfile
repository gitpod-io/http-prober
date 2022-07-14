FROM golang:1.18.3 as builder

WORKDIR /workspace

COPY ./*.go ./
COPY ./go.mod ./
COPY ./go.sum ./
COPY ./Makefile Makefile
RUN go mod download
RUN make build


FROM alpine:3.16.0
WORKDIR /
COPY --from=builder /workspace/http-prober /usr/bin/http-prober

ENTRYPOINT ["/usr/bin/http-prober"]