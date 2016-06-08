# RootFS

This directory stores all files that should be copied to the rootfs of a
Docker container. The files should be stored according to the correct
directory structure of the destination container. 


## Dockerfile

A Dockerfile in the rootfs is used to build the image. Where possible,
compilation should not be done in this Dockerfile, since we are
interested in deploying the smallest possible images.

