import docker
import logging
import os
import tarfile
import requests
from swiftclient.service import SwiftService, SwiftError
import time
import json
import subprocess

DEBUG = os.environ.get('DEBUG') in ('true', '1')
logging.basicConfig(level=logging.ERROR)
logging.getLogger("requests").setLevel(logging.CRITICAL)
logging.getLogger("swiftclient").setLevel(logging.CRITICAL)
logger = logging.getLogger(__name__)


def log_output(stream, decode):
    error = False
    for chunk in stream:
        if 'error' in chunk:
            error = True
            print(chunk.decode('utf-8'))
        elif decode:
            stream_chunk = chunk.get('stream')
            if stream_chunk:
                print(stream_chunk.replace('\n', ''))
        elif DEBUG:
            print(chunk.decode('utf-8'))
    if error:
        exit(1)


def download_file():
    if os.getenv('BUILDER_STORAGE') == "swift":
        with SwiftService(options={"auth_version": '3',
                          "os_project_id": os_project_id,
                          "os_user_id": os_user_id,
                          "os_password": os_password,
                          "os_auth_url": os_auth_url}) as swift:
            try:
                print('listing object store containers')
                # list_options = {"prefix": "archive_2016-01-01/"}
                list_parts_gen = swift.list(container=space_uuid)
                print(list_parts_gen)
                for page in list_parts_gen:
                    print(page)
                    if page["success"]:
                        for obj in page["listing"]:
                            if obj["name"].startswith(app_guid):
                                objects = [
                                    obj["name"] for obj in page["listing"]
                                ]
                        print objects

                        for down_res in swift.download(
                                container=space_uuid,
                                objects=objects,
                                options={
                                    'out_file': tar_path
                                }):
                            if down_res['success']:
                                logger.info("'%s' downloaded" % down_res['object'])
                            else:
                                logger.error("'%s' download failed" % down_res['object'])
                                logger.debug(down_res)
                    else:
                        raise page["error"]
            except SwiftError as e:
                logger.error(e.value)

tar_path = os.getenv('TAR_PATH', '/app/droplet.tgz')
space_uuid = os.getenv('SPACE_UUID', 'test1')
#app_guid = os.getenv('APP_GUID', '')
app_guid = ''

# TODO: not pass these as env vars.  store them on the kubernetes machines and use volume mount.
os_project_id = os.getenv('OS_PROJECT_ID', '43b6c705b03645a1b7aca19f0426bf58')
os_user_id = os.getenv('OS_USER_ID', 'e94915d69cfa43929f43185aa800a75c')
os_password = os.getenv('OS_PASSWORD', 'xxxx')
os_auth_url = os.getenv('OS_AUTH_URL', 'https://identity.open.softlayer.com/v3')

kubernetes_registry_host = os.getenv('K8S_REGISTRY_HOST', 'localhost')
kubernetes_registry_port = os.getenv('K8S_REGISTRY_PORT', '5000')
image_name = os.getenv('IMG_NAME', "first-droplet")
kubernetes_server = os.getenv('K8S_SERVER', '9.37.192.140')
# this should be the space_uuid
namespace = 'linsun'

logger.info("tar_path is %s", tar_path)
print tar_path

if tar_path:
    download_file()
logger.info("download tar file complete to %s", tar_path)
print "download tar file complete"

#with tarfile.open("apptar", "r:gz") as tar:
#    tar.extractall("/app/")
#log("extracting tar file complete")
client = docker.Client(version='auto')
registry = kubernetes_registry_host + ":" + kubernetes_registry_port
if ':' in image_name:
    image_name, image_tag = image_name.split(":", 1)
else:
    image_tag = None

repo = registry+'/' + namespace + '/' + image_name

build_response = [line for line in client.build(
    tag=repo, stream=True, decode=True, rm=True, path='/app'
)]

print build_response
print("Pushing to registry")
#time.sleep(30)

# TODO: create user's namespace if it doesn't exist yet

push_response = [line for line in client.push(registry+'/' + namespace + '/' + image_name, tag=image_tag, stream=True, insecure_registry=True)]
print push_response

deployment_name = image_name
# create the deployment
subprocess.call(['kubectl', '--server=http://' + kubernetes_server + ':8080/', 'run', deployment_name,
                 registry+'/' + namespace + '/' + image_name, '--port=8080', '--namespace=' + namespace])
