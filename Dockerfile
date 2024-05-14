# Use a builder image to build the binary
FROM golang:1.21 as builder
ARG VERSION
WORKDIR /app

# We copy everything: it's easier and avoids errors due to missing files.
COPY . .
RUN go mod download

# CGO_ENABLED=0 is **required** in order to build a static binary, otherwise
# it will fail to run in the Alpine image.
RUN make  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 VERSION=${VERSION} build

# ----
FROM alpine:3.19
ARG VERSION
WORKDIR /opt/majordomo

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/build/bin/majordomo-v${VERSION}-g*_linux-amd64 ./majordomo

# Change ownership of the binary to our non-root user
#RUN groupadd majordomo && useradd -g majordomo majordomo
RUN addgroup -S majordomo && adduser -S majordomo -G majordomo

RUN chown majordomo:majordomo majordomo
USER majordomo

EXPOSE 5000

# Command to run the executable
CMD ["./majordomo", "--port", "5000", "--config", "/etc/majordomo/config.yaml"]
