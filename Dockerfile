FROM golang:1.24.0-alpine3.21
WORKDIR /app
COPY ./src/go.mod ./src/go.sum ./
RUN go mod download
COPY ./src/ ./
RUN go build -o /netpartctrl
CMD [ "/netpartctrl" ]
