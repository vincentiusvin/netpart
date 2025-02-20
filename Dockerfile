FROM golang:1.24.0-alpine3.21 AS build
WORKDIR /app
COPY ./go.* .
RUN go mod download
COPY . .
RUN go build -o /netpartctrl


FROM docker:28.0.0-rc.3-dind-alpine3.21
COPY --chmod=744 ./entrypoint.sh ./entrypoint.sh
COPY --from=build /netpartctrl /netpartctrl
ENTRYPOINT ./entrypoint.sh
