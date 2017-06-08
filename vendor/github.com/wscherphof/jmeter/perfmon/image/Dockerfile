FROM busybox

COPY ServerAgent-2.2.1.tar.gz /src/

VOLUME /dst

ENTRYPOINT ["tar", "-xzf", "/src/ServerAgent-2.2.1.tar.gz", "-C", "/dst"]
