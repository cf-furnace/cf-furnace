FROM cloudfoundry/cflinuxfs2
ADD droplet /home/vcap
ADD launcher /tmp/lifecycle/launcher
RUN chown -R vcap:vcap /home/vcap /tmp/lifecycle
EXPOSE 8080
CMD ["/tmp/lifecycle/launcher", "app", "node server.js", ""]
