# Build Phase------------------------------------------------
FROM golang:1.18.0-alpine3.15 as build 

# Setup working dir
RUN mkdir -p /firebase_go_auth
WORKDIR /firebase_go_auth

# Copy and add files to workdir
ADD api /firebase_go_auth/api/
ADD email /firebase_go_auth/email/
ADD firebase_conn /firebase_go_auth/firebase_conn/
ADD utils /firebase_go_auth/utils/

COPY main.go /firebase_go_auth/
COPY go.sum /firebase_go_auth/
COPY go.mod /firebase_go_auth/

COPY .env /firebase_go_auth/
COPY serviceAccountKey.json /firebase_go_auth/

# Build init
RUN CGO_ENABLED=0 go build -o /firebase_go_auth/

# Image Serve Phase-------------------------------------------
FROM alpine:3.14

COPY --from=build /firebase_go_auth /firebase_go_auth

# Expose Ports
EXPOSE 8080

# Start Server
CMD cd firebase_go_auth && ./firebase_go_auth | tee firebase_go_auth.log
