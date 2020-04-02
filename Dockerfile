FROM centos
ADD prometheus-cert-exporter /tmp/
EXPOSE 8080
ENTRYPOINT /tmp/prometheus-cert-exporter