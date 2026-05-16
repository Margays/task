FROM golang:1.26.3 as base

FROM base AS workspace
ARG REMOTE_USER
ARG REMOTE_UID
ARG REMOTE_GID
RUN addgroup --gid ${REMOTE_GID} ${REMOTE_USER}
RUN adduser --disabled-password --shell /bin/bash --uid ${REMOTE_UID} --gid ${REMOTE_GID} ${REMOTE_USER}
ENV HOME /home/${REMOTE_USER}
ENV LC_ALL=C.UTF-8
USER ${REMOTE_USER}
