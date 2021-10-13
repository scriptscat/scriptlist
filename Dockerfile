FROM busybox:1.33.1-glibc

WORKDIR /www

COPY scriptlist .

RUN chmod +x scriptlist

ENTRYPOINT ["./scriptlist"]
