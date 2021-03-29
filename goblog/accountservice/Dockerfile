FROM golang AS builder

WORKDIR /go/src/app
COPY . /go/src/app
RUN CGO_ENABLED=0 go build -o /bin/accountservice .

FROM gcr.io/distroless/base
EXPOSE 6767
COPY --from=builder /bin/accountservice /accountservice
ENTRYPOINT ["/accountservice"]
