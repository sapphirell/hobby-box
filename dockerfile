FROM golang:1.18.2
#设置代理
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
#运行目录
RUN mkdir /app
ADD . /app
WORKDIR /app
#更换环境变量为正式文件
RUN rm -f .env && cp .env.prod .env
#编译
RUN go build -o main .
#运行
CMD ["./main"]