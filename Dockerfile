FROM golang:1.9.4-alpine3.7

RUN apk --update add \
bash \
curl \
gcc \
gettext \
git \
jq \
less \
libintl \
make \
openssh \
nodejs \
sudo \
vim \
wget \
zsh

RUN npm install serverless -g
RUN wget https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 \
&& chmod a+x dep-linux-amd64 \
&& mv dep-linux-amd64 /usr/bin/dep

RUN adduser -D developer
RUN echo "developer ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers
USER developer
WORKDIR /home/developer
RUN wget https://github.com/robbyrussell/oh-my-zsh/raw/master/tools/install.sh -O - | zsh || true
ENV SHELL=/bin/zsh
ENV PATH=$PATH:/home/developer/.local/bin

WORKDIR /go/src

ENTRYPOINT ["/bin/zsh", "-c", "while sleep 3600; do :; done"]
