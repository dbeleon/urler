FROM tarantool/tarantool:2.11.0

USER root

COPY urler/internal/urls-tnt /opt/tarantool

RUN ls -al /opt/tarantool

RUN rm -rf /opt/tarantool/data

ENTRYPOINT ["tarantool", "/opt/tarantool/init.lua"]