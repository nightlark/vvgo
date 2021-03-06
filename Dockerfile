FROM node:13.12 as node
COPY package.json .
COPY yarn.lock .
RUN yarn install

FROM golang:1.14 as builder

ARG GITHUB_REF
ARG GITHUB_SHA

ENV CGO_ENABLED=0 GOOS=linux GO111MODULE=on

WORKDIR /go/src/github.com/virtual-vgo/vvgo
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN BIN_PATH=/ make vvgo

FROM builder as tester
CMD ["make", "test"]

FROM alpine:3.4 as vvgo
RUN apk add --no-cache ca-certificates apache2-utils
COPY ./public /public
COPY --from=builder vvgo /vvgo
COPY --from=node node_modules /public/npm
EXPOSE 8080
CMD ["/vvgo"]
ENTRYPOINT ["/vvgo"]
