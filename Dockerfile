ARG BASE_IMAGE=alpine
ARG USER=service
ARG WORKDIR=/home/${USER}/app
ARG OUT_FILE=out/service

FROM ${BASE_IMAGE} as prep

ARG BUILD_DEPENDENCIES="make gcc go musl-dev file"
ARG WORKDIR
ARG USER

RUN apk update --no-cache \
  && apk upgrade \
  && apk add ${BUILD_DEPENDENCIES} \
  && adduser -D ${USER}

WORKDIR ${WORKDIR}
USER ${USER}

COPY --chown=${USER}:${USER} . .


FROM prep as build

ARG OUT_FILE
ARG CGO_ENABLED=0

RUN make image-executable \
  && ldd ${OUT_FILE} \
  && file ${OUT_FILE}


FROM scratch

ARG WORKDIR

COPY --from=build ${WORKDIR}/out .

CMD [ "/pwa-server" ]
