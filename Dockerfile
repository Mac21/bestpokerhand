FROM golang:latest as build

WORKDIR /app
COPY . .

RUN go mod download
RUN go install -v -tags=jsoniter ./...

FROM gcr.io/distroless/base-debian12
ENV HOST="0.0.0.0"
ENV PORT="5000"

COPY --from=build /app/templates /templates
COPY --from=build /app/static /static
COPY --from=build /go/bin/app /
CMD ["/app"]
