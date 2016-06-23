FROM cloudfoundry/cflinuxfs2
ADD cache stager /
ADD docker /usr/bin
ADD launcher.dockerfile /tmp/docker/Dockerfile
RUN chmod -R 1777 /tmp
CMD /stager
