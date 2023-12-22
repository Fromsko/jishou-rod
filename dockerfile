# 设置基础镜像为 ainow/alpine-chrome
FROM ainow/alpine-chrome

# 设置东八区，北京时间
ENV TZ=Asia/Shanghai

# 设置工作目录
WORKDIR /app

# 定义容器启动时的默认命令
CMD [ "./server" ]
