FROM golang:1.18-alpine

WORKDIR /usr/app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY ./ ./

RUN export GOBIN=$PWD
RUN go build -o /k8s-client

# RUN go install -o /k8s-client
EXPOSE 5050

CMD [ "/k8s-client" ]





