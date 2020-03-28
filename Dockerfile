FROM golang:1.14-alpine as build-env

WORKDIR /workspace
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY main.go .
COPY helpers.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/kaas

FROM alpine:3
COPY --from=build-env /go/bin/kaas /usr/bin/kaas

RUN apk add --no-cache ca-certificates curl git openssh terraform wget

# Download the typhoon ct provider
RUN wget https://github.com/poseidon/terraform-provider-ct/releases/download/v0.4.0/terraform-provider-ct-v0.4.0-linux-amd64.tar.gz && \
    tar xzf terraform-provider-ct-v0.4.0-linux-amd64.tar.gz && \
    mkdir -p  ~/.terraform.d/plugins/ && \
    mv terraform-provider-ct-v0.4.0-linux-amd64/terraform-provider-ct ~/.terraform.d/plugins/terraform-provider-ct_v0.4.0 && \
    rm -r terraform-provider-ct*

COPY tf/ tf/

ENTRYPOINT ["/usr/bin/kaas"]