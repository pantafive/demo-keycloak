FROM ubuntu:20.04 AS builder
ENV DEBIAN_FRONTEND=noninteractive

ENV PYTHON_VERSION=3.10.4
ENV PATH=/root/.poetry/bin:$PATH

RUN apt-get update
RUN apt-get install -y wget build-essential libreadline-gplv2-dev libncursesw5-dev \
 libssl-dev libsqlite3-dev tk-dev libgdbm-dev libc6-dev libbz2-dev libffi-dev zlib1g-dev

ADD https://www.python.org/ftp/python/$PYTHON_VERSION/Python-$PYTHON_VERSION.tgz ./
RUN tar -xzf Python-$PYTHON_VERSION.tgz && rm Python-$PYTHON_VERSION.tgz

RUN cd Python-$PYTHON_VERSION && ./configure --enable-optimizations && make install

ADD https://raw.githubusercontent.com/python-poetry/poetry/master/get-poetry.py get-poetry.py
RUN python3 get-poetry.py && rm get-poetry.py

COPY pyproject.toml poetry.lock ./
RUN poetry config virtualenvs.create false \
 && poetry install --no-dev --no-root --no-interaction --no-ansi


ENV PATH=$PATH:/usr/local/go/bin

WORKDIR /tmp
ADD https://go.dev/dl/go1.18.linux-amd64.tar.gz ./
RUN tar -C /usr/local -xzf go1.18.linux-amd64.tar.gz

WORKDIR /build-go
COPY go.mod go.sum ./
RUN go mod download
COPY backend_go/ ./
RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 go build -ldflags "-s -w" -o /backend_go .


FROM python:3.10-slim AS backend
COPY --from=builder /usr/local/lib/python3.10/site-packages /usr/local/lib/python3.10/site-packages
COPY --from=builder /usr/local/bin/uvicorn /usr/local/bin/uvicorn
WORKDIR /app
COPY backend_py ./
USER nobody
CMD [ "sh", "-c", "uvicorn main:app --host 0.0.0.0 --port 3000"]


FROM alpine:3 as backend_go
COPY --from=builder /backend_go /
CMD ["/backend_go"]


FROM node:alpine as frontend
COPY frontend/package*.json .
RUN npm install
COPY frontend/app.js frontend/index.html ./
CMD [ "npm", "start" ]
