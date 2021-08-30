# CONTEXT ../..
FROM golang
ADD . /build
WORKDIR /build/cmd/server
ENV CGO_ENABLED=0
RUN go build -o /result/server .

FROM scratch
COPY --from=build /result/server /app/bin/server
EXPOSE 8080
ENTRYPOINT ["/app/bin/server", "-config=/app/etc/config.yaml"]
