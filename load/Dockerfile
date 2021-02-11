FROM golang:alpine
COPY load.go /
COPY go.mod /
#COPY go.sum /
RUN go build -o /load /load.go

# Copy the exe into a smaller base image
FROM alpine
COPY --from=0 /load /load
CMD /load
