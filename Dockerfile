FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

FROM golang:1.25-alpine AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist
RUN go build -o easy-qfnu-api-lite .

FROM alpine:3.22
WORKDIR /app
COPY --from=backend-builder /app/easy-qfnu-api-lite ./
EXPOSE 8141
ENV PORT=8141
ENV GIN_MODE=release
CMD ["./easy-qfnu-api-lite"]
