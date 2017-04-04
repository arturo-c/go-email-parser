FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/arturo-c/go-email-parser

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
ENV GLIDE_VERSION 0.11.0
ENV GLIDE_DOWNLOAD_URL https://github.com/Masterminds/glide/releases/download/v${GLIDE_VERSION}/glide-v${GLIDE_VERSION}-linux-amd64.tar.gz

RUN curl -fsSL "$GLIDE_DOWNLOAD_URL" -o glide.tar.gz \
    && tar -xzf glide.tar.gz \
    && mv linux-amd64/glide /usr/bin/ \
    && rm -r linux-amd64 \
    && rm glide.tar.gz

WORKDIR $GOPATH/src/github.com/arturo-c/go-email-parser

RUN glide install
RUN go install github.com/arturo-c/go-email-parser
RUN go build

# Document that the service listens on port 8080.
EXPOSE 8080

CMD ./go-email-parser
