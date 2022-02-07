ARG BASE_IMAGE=alpine
ARG USER=service
ARG WORKDIR=/home/${USER}

FROM ${BASE_IMAGE} as prep

ARG RUNTIME_DEPENDENCIES="git-daemon"
ARG WORKDIR
ARG USER

RUN apk update --no-cache \
  && apk upgrade \
  && apk add ${RUNTIME_DEPENDENCIES} \
  && adduser -D ${USER}

USER ${USER}

WORKDIR ${WORKDIR}
