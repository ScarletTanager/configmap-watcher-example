FROM golang:alpine
COPY kn-watch.go /
COPY go.mod /
COPY go.sum /
RUN go build -o /knwatch /kn-watch.go

# Copy the exe into a smaller base image
FROM alpine
COPY --from=0 /knwatch /knwatch
CMD /knwatch
