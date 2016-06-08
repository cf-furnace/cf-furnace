FROM cloudfoundry/cflinuxfs2:latest
COPY droplet.tgz /app/
COPY start.sh /app/
RUN chown vcap:vcap /app && cd /app && su vcap -c "tar zxf droplet.tgz" && rm droplet.tgz
EXPOSE 8080
USER vcap
WORKDIR /app
ENV HOME /app/app
ENV PORT 8080
ENV TMPDIR /app/tmp
CMD ["/bin/bash", "/app/start.sh", "/app", "npm start"]

