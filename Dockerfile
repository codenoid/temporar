FROM golang:alpine AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy the code into the container and download dependency using go mod
COPY . .

RUN apk --no-cache add ca-certificates

# Build Svelte static files
RUN apk add --update nodejs npm
WORKDIR /build/frontend
RUN npm install
RUN npm run build
WORKDIR /build


# Build the Go Binary
RUN go mod tidy
RUN go build -trimpath -o main .

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/main .

############################
# STEP 2 build a small image
############################
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /dist/main /

# Command to run the executable
ENTRYPOINT ["/main"]
