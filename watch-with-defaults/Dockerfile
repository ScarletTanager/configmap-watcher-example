FROM golang:alpine
COPY watch-with-defaults.go /
COPY go.mod /
COPY go.sum /
RUN go build -o /watch-with-defaults /watch-with-defaults.go

# Copy the exe into a smaller base image
FROM alpine
COPY --from=0 /watch-with-defaults /watch-with-defaults
CMD /watch-with-defaults
