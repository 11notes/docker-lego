# ╔═════════════════════════════════════════════════════╗
# ║                       SETUP                         ║
# ╚═════════════════════════════════════════════════════╝

  # arguments
  ARG APP_UID=1000 \
      APP_GID=1000

  # foreign image layers
  FROM 11notes/util AS util
  FROM 11notes/distroless:lego AS distroless-lego

# ╔═════════════════════════════════════════════════════╗
# ║                       BUILD                         ║
# ╚═════════════════════════════════════════════════════╝

  # build distroless lego-cron
  FROM golang:1.24-alpine AS build

  ARG APP_ROOT

  ENV BUILD_ROOT=/go/lego-cron \
      BUILD_BIN=/go/lego-cron/lego-cron \
      CGO_ENABLED=0

  COPY --from=util /usr/local/bin/ /usr/local/bin
  COPY ./go/lego-cron /go/lego-cron

  USER root

  RUN set -ex; \
    apk --update --no-cache add \
      build-base \
      upx;

  RUN set -ex; \
    cd ${BUILD_ROOT}; \
    go mod tidy;

  RUN set -ex; \
    cd ${BUILD_ROOT}; \
    go build -ldflags="-extldflags=-static" -o ${BUILD_BIN} main.go;

  RUN set -ex; \
    eleven checkStatic ${BUILD_BIN}; \
    eleven strip ${BUILD_BIN}; \
    mkdir -p /distroless/usr/local/bin; \
    cp ${BUILD_BIN} /distroless/usr/local/bin;

  RUN set -ex; \
    mkdir -p /distroless${APP_ROOT}/etc; \
    mkdir -p /distroless${APP_ROOT}/var;

# ╔═════════════════════════════════════════════════════╗
# ║                       IMAGE                         ║
# ╚═════════════════════════════════════════════════════╝

  # :: HEADER
    FROM scratch

    ARG APP_ROOT \
        APP_UID \
        APP_GID

    ENV APP_ROOT=${APP_ROOT}

    COPY --from=distroless-lego --chown=${APP_UID}:${APP_GID} / /
    COPY --from=build --chown=${APP_UID}:${APP_GID} /distroless/ /

  # :: PERSISTENT DATA
    VOLUME ["${APP_ROOT}/etc", "${APP_ROOT}/var"]

  # :: START
    USER ${APP_UID}:${APP_GID}
    ENTRYPOINT ["/usr/local/bin/lego-cron"]