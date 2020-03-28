FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY goChat /goChat
COPY html /html
EXPOSE 8080
ENTRYPOINT /goChat
LABEL Name=goChat Version=0.0.1
