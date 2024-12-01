FROM golang:1.23.3

# Set working directory
WORKDIR /app
USER root

# Cache this step
COPY go.mod go.sum ./
RUN go mod download

# Copy source files, from . to .
COPY . .

# Ensure entrypoint is executable
RUN chmod +x ./deployment/docker/entrypoint.sh

ENTRYPOINT ["./deployment/docker/entrypoint.sh"]