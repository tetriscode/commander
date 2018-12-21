FROM tetriscodes/alpine-go-kafka:1.9.3
MAINTAINER "Wes Richardet<wes@tetriscodes.com>"
# Copy the local package files to the container's workspace.
ADD ./ /go/src/github.com/tetriscode/commander

WORKDIR /go/src/github.com/tetriscode/commander
RUN go install github.com/tetriscode/commander

ENTRYPOINT /go/bin/commander
EXPOSE 8083
