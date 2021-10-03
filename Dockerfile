
FROM wehack/alpine-glibc
WORKDIR /app

COPY vsat64 .
COPY docker-entrypoint.sh .
RUN chmod +x ./docker-entrypoint.sh
EXPOSE 1025/tcp

# install nmap
# RUN apk add nmap nmap-scripts

ENTRYPOINT ["./docker-entrypoint.sh"]
CMD ["./vsat64"]