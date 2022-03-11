ARG BASE_IMAGE=alpine
ARG USER=builder
ARG WORKDIR=/home/${USER}
ARG OUT_DIR="${WORKDIR}/out"
ARG OUT_FILE="${OUT_DIR}/bin/git"
ARG REF=master

FROM ${BASE_IMAGE} as prep

ARG BUILD_DEPENDENCIES="make autoconf linux-headers musl-dev gcc zlib-dev zlib-static file"
ARG WORKDIR
ARG USER
ARG REF

RUN apk update --no-cache \
  && apk upgrade \
  && apk add ${BUILD_DEPENDENCIES} \
  && adduser -D ${USER}

WORKDIR "${WORKDIR}"
USER ${USER}

ADD --chown=${USER}:${USER} https://github.com/git/git/archive/refs/heads/${REF}.zip .

RUN unzip "${REF}"
WORKDIR "${WORKDIR}/git-${REF}"

FROM prep as configure

ARG OUT_DIR
ARG CFLAGS="-static"
ARG NO_TCLTK=true

RUN make configure && ./configure --prefix="${OUT_DIR}" --with-ssl --with-curl


FROM configure as build

RUN make -j6


FROM build as install

RUN make -j6 install


FROM scratch

ARG WORKDIR
ARG OUT_FILE

COPY --from=install ${OUT_FILE} .
