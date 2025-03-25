FROM alpine:3.13.2 AS builder

# Install all dependencies required for compiling busybox
RUN apk add gcc musl-dev make perl linux-headers build-base git

RUN git clone git://busybox.net/busybox.git --branch master --single-branch

WORKDIR /busybox

# Copy the busybox build config (limited to httpd)
COPY ./tests/.config .config

# Compile and install busybox
RUN make -j && make install -j

# Create a non-root user to own the files and run our server
RUN adduser -D static

# Switch to the scratch image
FROM scratch

# Copy over the user
COPY --from=builder /etc/passwd /etc/passwd

# Copy the busybox static binary
COPY --from=builder /busybox/_install/bin/busybox /

# Uploads a blank default httpd.conf
# This is only needed in order to set the `-c` argument in this base file
# and save the developer the need to override the CMD line in case they ever
# want to use a httpd.conf
COPY ./tests/httpd.conf /httpd.conf

WORKDIR /home/static

# Run busybox httpd
CMD ["/busybox", "httpd", "-p", "3000", "-vv", "-f", "httpd.conf", "/home/static"]

EXPOSE 3000

# Use our non-root user
USER static