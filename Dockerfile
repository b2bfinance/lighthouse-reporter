FROM golang:1.12-buster as lhreporter

COPY . /go/src/github.com/b2bfinance/lighthouse-reporter

RUN go install github.com/b2bfinance/lighthouse-reporter/cmd/lhreporter

FROM node:12-buster

RUN apt-get update --fix-missing && apt-get -y upgrade
RUN apt-get install -y sudo xvfb dbus-x11 --no-install-recommends

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

COPY --from=lhreporter /go/bin/lhreporter /usr/local/bin/lhreporter
COPY entrypoint.sh /entrypoint.sh
COPY start-headless-chrome.sh /start-headless-chrome.sh

ENTRYPOINT ["/entrypoint.sh"]
CMD []
