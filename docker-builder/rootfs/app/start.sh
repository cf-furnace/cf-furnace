echo "start.sh starts"
if [ -z "$1" ]; then
  echo "usage: $0 <app dir> <command to run>" >&2
  exit 1
fi

cd "$1/app"

if [ -d .profile.d ]; then
  for env_file in .profile.d/*; do
    echo "sourcing " $env_file
    source $env_file
  done
fi

shift

eval "$@"
