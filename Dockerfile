FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -trimpath -o /out/alerter ./cmd/alerter

FROM scratch
COPY --from=build /out/alerter /alerter
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 3000
USER 65534:65534
ENTRYPOINT ["/alerter"]
