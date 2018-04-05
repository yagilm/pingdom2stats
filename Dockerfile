FROM scratch
ADD ca-certificates.crt /etc/ssl/certs/
ADD pingdom2stats-docker /
ENTRYPOINT ["/pingdom2stats-docker"]
CMD ["--help"]
