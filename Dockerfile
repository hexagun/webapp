# # Stage 1: Build React app
# FROM node:22-alpine as npm-builder

# WORKDIR /app
# COPY ./frontend .
# RUN npm install
# RUN npm run build
# RUN ls -l
# RUN ls -l ./dist

# # Stage 2: Build Go server
# FROM golang:1.22-alpine as go-builder

# WORKDIR /app

# ADD server.go /app/
# ADD go.mod /app/
# ADD go.sum /app/

# RUN go mod download
# RUN CGO_ENABLED=0 GOOS=linux go build -o main
# RUN ls -l

# # Stage 3: Run the Go server
# FROM alpine:latest

# WORKDIR /app
# RUN mkdir /frontend
# COPY --from=npm-builder /app/dist /app/frontend/dist
# COPY --from=go-builder /app/main /app/main

# CMD ["go", "run", "main.go"]

FROM node:22-alpine as npm-builder

WORKDIR /app
COPY ./frontend .
RUN npm install
RUN npm run build
RUN ls -l
RUN ls -l ./dist

FROM golang:1.22-alpine

WORKDIR /app

RUN mkdir /frontend
COPY --from=npm-builder /app/dist /app/frontend/dist

COPY go.mod ./
RUN go mod tidy

COPY . .
RUN go build -o server .

CMD ["./server"]