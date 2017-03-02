FROM alpine

MAINTAINER Zhang Peihao <zhangpeihao@gmail.com>

ADD gocmd.linux /app/tools/gocmd/gocmd
ADD gocmd.run.sh /app/tools/gocmd/gocmd.run.sh
ADD gocmd-scripts /app/tools/gocmd/gocmd-scripts

WORKDIR /app/tools/gocmd

ENTRYPOINT ["/app/tools/gocmd/gocmd.run.sh"]

EXPOSE 8001
