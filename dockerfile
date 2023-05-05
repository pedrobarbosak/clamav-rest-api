##############################
#   Build Container
FROM golang:alpine as build

ENV CGO_ENABLED=0

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o clamav-rest-api .

##############################
#   Run Container
FROM alpine:latest as run

WORKDIR /app

# Install clamav
RUN apk --no-cache add clamav clamav-libunrar

# Setup required folders and permissions
RUN mkdir /run/clamav && chown clamav:clamav /run/clamav

# Setup clamd
RUN sed -i 's/^#Foreground .*$/Foreground true/g' /etc/clamav/clamd.conf
RUN sed -i 's/^#TCPSocket .*$/TCPSocket 3310/g' /etc/clamav/clamd.conf
RUN sed -i 's/^#Foreground .*$/Foreground true/g' /etc/clamav/freshclam.conf

# Download databases
RUN freshclam --quiet

# Copy necessary files from Build Container
COPY --from=build /app/clamav-rest-api /usr/bin/.

# Copy entry point
COPY entrypoint.sh /usr/bin/.

# Set runnable permissions
RUN chmod +x /usr/bin/entrypoint.sh

# Env variables
ENV PORT=8080

# Run
ENTRYPOINT ["entrypoint.sh"]