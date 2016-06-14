FROM cloudfoundry/cflinuxfs2
ADD cache stager install_docker /
ADD launcher.dockerfile /tmp/docker/Dockerfile
RUN chmod -R 1777 /tmp
RUN /install_docker
CMD /stager
