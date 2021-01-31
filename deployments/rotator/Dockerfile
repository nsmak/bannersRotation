FROM golang:1.15.7 as builder

ENV BIN_FILE /opt/rotator/rotator-app
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/rotator/*

FROM alpine:3.9

ENV BIN_FILE "/opt/rotator/rotator-app"
COPY --from=builder ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/rotator/config.json
COPY ./configs/rotator.json ${CONFIG_FILE}

CMD ${BIN_FILE} -config ${CONFIG_FILE}