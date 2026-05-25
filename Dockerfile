FROM golang:1.25-alpine

WORKDIR /app

COPY . ./
RUN go env -w GOPROXY=https://goproxy.cn,direct \
    && go mod download

RUN go build -o /product-mall

EXPOSE 3000

CMD [ "/product-mall" ]