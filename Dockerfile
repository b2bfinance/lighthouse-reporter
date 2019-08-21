FROM golang:1.12-buster as lhreporter

COPY . /go/src/github.com/b2bfinance/lighthouse-reporter

RUN go install github.com/b2bfinance/lighthouse-reporter/cmd/lhreporter

FROM node:12-buster

RUN apt-get update --fix-missing && apt-get -y upgrade

RUN wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add - \
    && sh -c 'echo "deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google.list' \
    && apt-get update \
    && apt-get install -y google-chrome-unstable --no-install-recommends \
    && rm -rf /var/lib/apt/lists/* \
    && rm -rf /src/*.deb

RUN groupadd --system chrome && \
    useradd --system --create-home --gid chrome --groups audio,video chrome && \
    mkdir --parents /home/chrome/reports && \
    chown --recursive chrome:chrome /home/chrome

RUN npm i lighthouse -g

COPY package.json .
RUN npm i --production

USER chrome

COPY --from=lhreporter /go/bin/lhreporter /usr/local/bin/lhreporter

ENTRYPOINT /usr/local/bin/lhreporter
