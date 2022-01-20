ARG BASE_IMAGE=alpine
ARG USER=service
ARG WORKDIR=/home/${USER}/app
ARG OUT_FILE=out/service

FROM ${BASE_IMAGE} as prep

ARG WORKDIR
ARG USER
ARG BUILD_DEPENDENCIES="make go gcc musl-dev file build-base"

RUN apk update --no-cache \
    && apk upgrade \
    && apk add ${BUILD_DEPENDENCIES} \
    && adduser -D ${USER}

USER ${USER}

COPY --chown=${USER}:${USER} . ${WORKDIR}


FROM prep as build

ARG WORKDIR
ARG OUT_FILE
ARG CGO_ENABLED=0

WORKDIR ${WORKDIR}

RUN make executable certs \
    && ldd ${OUT_FILE} \
    && file ${OUT_FILE}


FROM scratch

ARG WORKDIR

COPY --from=build ${WORKDIR}/out .

CMD [ "/pwa-server" ]
