FROM golang:1.16-alpine

WORKDIR /usr/app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY ./ ./

RUN export GOBIN=$PWD
RUN go install

# RUN go install -o /k8s-client

CMD [ "" ]





