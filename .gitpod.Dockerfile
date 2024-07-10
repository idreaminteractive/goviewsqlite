FROM gitpod/workspace-base:2024-06-26-08-49-45
# sqlite, litefs
RUN sudo apt-get update && sudo apt-get install -y sqlite

# before other go things
ENV GO_VERSION=1.22.3
# For ref, see: https://github.com/gitpod-io/workspace-images/blob/61df77aad71689504112e1087bb7e26d45a43d10/chunks/lang-go/Dockerfile#L10
ENV GOPATH=$HOME/go-packages
ENV GOROOT=$HOME/go
ENV PATH=$GOROOT/bin:$GOPATH/bin:$PATH
RUN curl -fsSL https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz | tar xzs \
    && printf '%s\n' 'export GOPATH=/workspace/go' \
                      'export PATH=$GOPATH/bin:$PATH' > $HOME/.bashrc.d/300-go


# install air
RUN go install github.com/air-verse/air@latest

# install task
RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b ~/.local/bin

RUN sudo apt -y install python3-pip
RUN pip install -U litecli

RUN go install github.com/a-h/templ/cmd/templ@v0.2.747
RUN go install github.com/idreaminteractive/goreload/cmd/goreload@v0.0.3



RUN sudo mkdir /data
RUN sudo chown -R gitpod /data

# alias all the things
RUN echo 'alias home="cd ${GITPOD_REPO_ROOT}"' | tee -a ~/.bashrc ~/.zshrc
