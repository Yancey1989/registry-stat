FROM busybox
ADD ./registry-stat /registry-stat
ADD ./docker-entrypoint.sh /docker-entrypoint.sh
ENTRYPOINT ["/bin/sh", "/docker-entrypoint.sh"]
