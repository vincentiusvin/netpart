FROM golang:1.24.0-alpine3.21
WORKDIR /app
COPY ./go.* .
RUN go mod download
COPY . .
RUN go build -o /netpartctrl
CMD [ "/netpartctrl" ]
