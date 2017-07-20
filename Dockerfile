FROM scratch
ADD ca-certificates.crt /etc/ssl/certs/
ADD pingdom2mysql-docker /
ENTRYPOINT ["/pingdom2mysql-docker"]
CMD ["--help"]
