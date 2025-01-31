# build
FROM golang AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY src/*.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /calunite

# deploy
FROM alpine AS deploy

WORKDIR /

COPY --from=builder /calunite /calunite

# default configuration
ENV CFG_PATH=/config/config.yml
ENV CRON="@every 15m"
ENV PROD_ID=CalUnite
ENV CONTENT_DIR=/wwwdata
ENV FILE_NAVIGATION=false
ENV ADDR=0.0.0.0
ENV PORT=8080

EXPOSE $PORT

ENTRYPOINT ["/calunite"]
