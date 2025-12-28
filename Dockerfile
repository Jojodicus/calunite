# build
FROM golang AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY src/*.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /calunite

# deploy
FROM gcr.io/distroless/static-debian12 AS deploy

WORKDIR /

# default configuration
ENV CFG_PATH=/config/config.yml
ENV CRON="@every 15m"
ENV PROD_ID=CalUnite
ENV CONTENT_DIR=/wwwdata
ENV FILE_NAVIGATION=false
ENV DOT_PRIVATE=true
ENV ADDR=0.0.0.0
ENV PORT=8080

EXPOSE $PORT

COPY --from=builder /calunite /calunite

ENTRYPOINT ["/calunite"]
