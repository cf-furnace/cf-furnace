FROM cloudfoundry/cflinuxfs2:latest

COPY launcher /app/
COPY run.sh /app/
RUN chown vcap:vcap /app

USER vcap
WORKDIR /app

CMD ["/bin/bash", "run.sh", "app", "$STARTCMD", "$ENVVARS"]


