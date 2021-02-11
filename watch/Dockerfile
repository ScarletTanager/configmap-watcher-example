FROM golang:alpine
COPY watch.go /
COPY go.mod /
COPY go.sum /
RUN go build -o /watch /watch.go

# Copy the exe into a smaller base image
FROM alpine
COPY --from=0 /watch /watch
CMD /watch
