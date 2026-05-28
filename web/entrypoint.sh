#!/bin/sh
# Supervise busybox httpd + haproxy. Exit if either dies so tini reaps the
# container and the orchestrator restarts.
set -e

cleanup() {
    kill "$HTTPD_PID" "$HAPROXY_PID" 2>/dev/null || true
}
trap cleanup EXIT INT TERM

httpd -p 8081 -h /var/www/html -f &
HTTPD_PID=$!

haproxy -W -db -f /usr/local/etc/haproxy/haproxy.cfg &
HAPROXY_PID=$!

# Poll: exit as soon as either child dies (portable across ash/dash/bash).
while kill -0 "$HTTPD_PID" 2>/dev/null && kill -0 "$HAPROXY_PID" 2>/dev/null; do
    sleep 2
done

exit 1
