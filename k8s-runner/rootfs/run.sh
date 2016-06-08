echo "run.sh starts"
if [ -z "$1" ]; then
  echo "usage: $0 <app dir> <command to run> <env_vars>" >&2
  sleep 300
  exit 1
fi


ls /app

chown vcap:vcap /app

cd /app

cp /app/droplet/droplet.tgz /app

tar zxf droplet.tgz

./launcher $1 "$2" "$3"