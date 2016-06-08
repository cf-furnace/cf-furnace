pushd ../rootfs
docker build -t dockerbuilder .
pop

docker tag dockerbuilder linsun/dockerbuilder
docker push linsun/dockerbuilder
