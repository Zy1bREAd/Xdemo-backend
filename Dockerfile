FROM golang:1.23-alpine
MAINTAINER OceanWang
WORKDIR /app
ENV XDEMO_SYSTEM_MODE=container \
    XDEMO_SYSTEM_HOST=0.0.0.0 XDEMO_SYSTEM_PORT=7077 XDEMO_SYSTEM_MODE=container \
    XDEMO_DATABASE_HOST=10.0.20.5 XDEMO_DATABASE_PORT=2206  XDEMO_DATABASE_DBNAME=oceantest XDEMO_DATABASE_DBUSER=oceantestusr XDEMO_DATABASE_DBPASSWORD=oceanwangpwd \
    XDEMO_REDIS_ADDR=10.0.20.5 XDEMO_REDIS_PORT=6379 XDEMO_REDIS_DB=0 XDEMO_REDIS_PASSWORD=  XDEMO_REDIS_TLS=false \
    XDEMO_DOCKER_HOST=tcp://10.0.20.5:5732 XDEMO_DOCKER_VERSION=1.43 \
    XDEMO_K8S_MODE=2 XDEMO_K8S_KUBECONFIG= \
    XDEMO_QUEUE_PROVIDER=redis XDEMO_QUEUE_PROCESSER=3

# 拷贝当前git目录所有内容到/app下
COPY . .
# 正常应该由一个地方copy过来，而不是放在git上面
# RUN mkdir -p /root/.kube && cp secret/config /root/.kube/config && cp -rf secret/home /
COPY ./secret/config /root/.kube/config
COPY ./secret/home /
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download
RUN go build -o xdemoapp
# APP 访问端口
EXPOSE 7077
CMD ["/app/xdemoapp"]