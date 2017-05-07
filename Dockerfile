FROM golang:1.8
ENV DIST /go/src/github.com/Yancey1989/registry-stat
ADD ./docker-entrypoint.sh /docker-entrypoint.sh
COPY . $DIST
RUN cd $DIST && go get ./... && go get .
ENTRYPOINT ["/bin/sh", "/docker-entrypoint.sh"]
