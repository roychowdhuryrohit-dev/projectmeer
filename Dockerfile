FROM node:21-alpine AS builder-container-node
WORKDIR /app
ADD ./assets /app
RUN apk add --update --no-cache git
RUN npm install && npm run build

FROM golang:1.21-alpine AS builder-container-go
WORKDIR /app
ADD . /app
RUN apk add --update --no-cache git
RUN apk --update --no-cache add ca-certificates && \
    update-ca-certificates
RUN cd /app && mkdir -p /app/bin && \
    CGO_ENABLED=0 go build -o /app/bin/meer -tags netgo

FROM scratch
WORKDIR /app

COPY --from=builder-container-node /app/build /app/assets/build
COPY --from=builder-container-go /app/bin/meer /app/bin/meer
COPY --from=builder-container-go /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt


ENTRYPOINT [ "/app/bin/meer" ]