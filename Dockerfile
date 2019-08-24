FROM alpine
ARG version

LABEL com.plasmabank.version=${version} \
      com.dennybaa.app=helloapp

ENV PORT=8080 \
    MONGODB_CONNTIMEOUT="" \
    MONGODB_URI="" \
    MONGODB_DATABASE="" \
    APP_ENV="dev"

COPY ./dist/helloapp-linux-x86_64 /usr/local/bin/helloapp
CMD /usr/local/bin/helloapp
