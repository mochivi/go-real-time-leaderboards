FROM golang:1.23.3

# Set working directory
WORKDIR /app
USER root

# Copy source files, from . to .
COPY . .

RUN go mod download
RUN chmod +x ./deployment/docker/entrypoint.sh

ENTRYPOINT ["./deployment/docker/entrypoint.sh"]