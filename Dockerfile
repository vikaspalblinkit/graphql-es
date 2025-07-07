# ---------- Builder Stage ----------
FROM golang:1.23.4 AS builder

WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Install gqlgen (used during code generation)
RUN go install github.com/99designs/gqlgen@latest

# Generate GraphQL code (only if gqlgen.yml is present and gqlgen is used)
ENV PATH="/go/bin:${PATH}"
RUN gqlgen generate


# Build the Go application
RUN go build -o server .

# ---------- Runtime Stage ----------
FROM golang:1.23.4

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/server .

# Optional: Copy other required runtime files (like .env or configs), if they exist
# You can uncomment the next line if you have a .env file and it's required at runtime
# COPY --from=builder /app/.env .

# Set the default command to run the server binary
CMD ["./server"]
