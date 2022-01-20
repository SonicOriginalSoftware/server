ARG BASE_IMAGE=alpine

FROM ${BASE_IMAGE} as prep

ARG USER=service
ARG WORKDIR=/home/${USER}/app
ARG BUILD_DEPENDENCIES="make go gcc musl-dev"

RUN apk update --no-cache \
    && apk upgrade \
    && apk add ${BUILD_DEPENDENCIES} \
    && adduser -D ${USER}

USER ${USER}

COPY --chown=${USER}:${USER} . ${WORKDIR}

FROM prep as build

ARG CGO_ENABLED=1

WORKDIR ${WORKDIR}

RUN make executable certs
