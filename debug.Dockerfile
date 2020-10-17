#first stage - builder
FROM golang:1.15.3-buster AS dependencies

WORKDIR /build

# Install the file util
RUN apt-get update && apt-get install -y file && \
    go get github.com/go-delve/delve/cmd/dlv
# upx removed due to binary corruption
#libc6-compat make upx

# https://stackoverflow.com/questions/32232655/go-get-results-in-terminal-prompts-disabled-error-for-github-private-repo/44247040#44247040
#RUN git config --global url."https://fcc***03:@github.company.com/".insteadOf "https://github.company.com/"
#ENV GOPRIVATE github.company.com

# Resolve and build Go dependencies as Docker cache 
COPY go.mod /build/go.mod
COPY go.sum /build/go.sum
ENV GO111MODULE=on
RUN go mod download

FROM dependencies AS compilation

# Add the main
COPY main.go /build/main.go

##################
################## ATTENTION
##################
# Add the modules dirs to avoid errors 
# cannot load github.company.com/user/repo/api/git: no matching versions for query "latest"
COPY api /build/api
COPY cmd /build/cmd
COPY config /build/config
COPY util /build/util

ARG BIN_NAME
ARG BIN_VERSION

# Cross-compile all versions
ENV BIN_NAME ${BIN_NAME:-unknown}
ENV BUILD_VERSION ${BIN_VERSION:-0.1.0}
ENV PLATFORMS "darwin linux windows"
ENV ARCHS "amd64"
#ENV ARCHS "386 amd64"
# Debugging with delv options using GC_FLAGS -gcflags="all=-N -l"
# More at https://kupczynski.info/2020/05/17/remote-debug-go-code.html
# https://github.com/igor-kupczynski/remote-debug-example/blob/master/Dockerfile.debug#L9
ENV GC_FLAGS "all=-N -l"

# Build the module name
COPY .git/ /build/.git/

# Injecting version info into the golang build https://github.com/Ropes/go-linker-vars-example
# https://github.com/Ropes/go-linker-vars-example, https://stackoverflow.com/questions/11354518/application-auto-build-versioning/11355611#11355611
# https://medium.com/@joshroppo/setting-go-1-5-variables-at-compile-time-for-versioning-5b30a965d33e
RUN export export FULL_NAME_GIT=$(git -C /build/.git remote -v | grep fetch | awk '{print $2}' | sed -En "s/git@//p" | sed -En "s/.git//p" | sed -En "s/:/\//p") && \
    export export FULL_NAME_HTTP=$(git -C /build/.git remote -v | grep fetch | awk '{print $2}' | sed -En "s/https:\/\///p") && \
    export GO_MODULE_FULL_NAME=$(if [ ! -z "$FULL_NAME_GIT" ]; then echo "$FULL_NAME_GIT"; else echo "$FULL_NAME_HTTP"; fi) && \
    export BUILD_GIT_SHA=$(git rev-parse --short HEAD) && \
    for GOOS in ${PLATFORMS}; do for GOARCH in ${ARCHS}; do BINARY="${BIN_NAME}-$GOOS-$GOARCH"; if [ $GOOS = "windows" ]; then BINARY="${BINARY}.exe"; fi; export BUILD_TIME="$(date -u +"%Y-%m-%d_%H:%M:%S_GMT")"; echo "Cross-compiling $GO_MODULE_FULL_NAME@$BUILD_GIT_SHA as ${BINARY}@$BUILD_VERSION at $BUILD_TIME"; GO_MODULE_FULL_NAME=$GO_MODULE_FULL_NAME BUILD_GIT_SHA=$BUILD_GIT_SHA BUILD_VERSION=$BUILD_VERSION BUILD_TIME=$BUILD_TIME GO111MODULE=$GO111MODULE CGO_ENABLED=0 GOARCH=$GOARCH GOOS=$GOOS GOPRIVATE=$GOPRIVATE go build -o ${BINARY} -gcflags="${GC_FLAGS}" -a -ldflags "-X ${GO_MODULE_FULL_NAME}/config.VersionBuildGoModule=${GO_MODULE_FULL_NAME} -X ${GO_MODULE_FULL_NAME}/config.VersionBuildNumber=${BUILD_VERSION} -X ${GO_MODULE_FULL_NAME}/config.VersionBuildGitSha=${BUILD_GIT_SHA} -X ${GO_MODULE_FULL_NAME}/config.VersionBuildTime=${BUILD_TIME}"; ls -la "/build/${BINARY}"; file "/build/${BINARY}"; chmod +x "/build/${BINARY}"; if [ $GOOS = "linux" ]; then sh -c "/build/${BINARY}"; sh -c "/build/${BINARY} version"; fi; done; done

# Compress the binaries
# It is not working with errors like https://github.com/golang/go/issues/19625
#RUN upx --lzma /build/${BIN_NAME}*

# Build the main container (Linux Runtime)
FROM ubuntu:20.04
WORKDIR /root/

ARG BIN_NAME
ENV BIN_NAME ${BIN_NAME:-unknown}

# Copy the linux amd64 binary, based on the arg (or else all files are copied) inspect with https://github.com/wagoodman/dive
COPY --from=compilation /build/${BIN_NAME}* /usr/local/bin/

# Debugging with delve https://kupczynski.info/2020/05/17/remote-debug-go-code.html
COPY --from=dependencies /go/bin/dlv /

# Debugging information based on the CG_FLAGS with delve https://kupczynski.info/2020/05/17/remote-debug-go-code.html
ENV CG_FLAGS_PORT ${CG_FLAGS_PORT:-40000}

# Command to be provided to the debugger, as it is required (CLIs not Servers)
ARG DEBUG_COMMAND
ENV DEBUG_COMMAND ${DEBUG_COMMAND}

# Move the bin to /usr/local/bin and make the entrypoint to point to it passing the params
# https://stackoverflow.com/questions/33439230/how-to-write-commands-with-multiple-lines-in-dockerfile-while-preserving-the-new/33439625#33439625
RUN echo "#!/bin/sh" > /usr/local/bin/entrypoint.sh && \
    # https://stackoverflow.com/questions/32727594/how-to-pass-arguments-to-shell-script-through-docker-run/40312311#40312311
    echo "/dlv --listen=:${CG_FLAGS_PORT} --headless=true --api-version=2 --accept-multiclient exec /usr/local/bin/${BIN_NAME} -- ${DEBUG_COMMAND}" >> /usr/local/bin/entrypoint.sh && \
    chmod +x /usr/local/bin/entrypoint.sh && \
    ln -s /usr/local/bin/${BIN_NAME}-linux-amd64 /usr/local/bin/${BIN_NAME}
    # Delete all binarines that are not the linux one https://www.cyberciti.biz/faq/find-command-exclude-ignore-files/
    # find /bin -type f -name "cloner-*" ! -path "*linux*" | xargs rm -f && \

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
