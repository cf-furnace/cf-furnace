FROM phusion/baseimage

ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update -qq && \
    apt-get install -y -q --no-install-recommends \
    python2.7 python-pip build-essential python-dev

RUN pip install docker-py==1.6.0 python-keystoneclient==1.6.0 python-swiftclient==3.0.0

RUN pip freeze

RUN curl -O https://storage.googleapis.com/bin.kuar.io/linux/kubectl
RUN chmod +x kubectl
RUN sudo cp kubectl /usr/local/bin/kubectl

COPY . /
COPY app/ app/

ENTRYPOINT ["python", "/deploy.py"]
