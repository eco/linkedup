FROM golang

RUN apt-get update && \
  apt-get install -y awscli && \
  aws configure set aws_access_key_id bddashdb2webcwhcbwdv23eb2abcd && \
  aws configure set aws_secret_access_key wdbhwchbwvbwb23123cebdchxyz && \
  aws configure set default.region us-west-1

WORKDIR /src
COPY . /src/
RUN make install

CMD sh scripts/start-all.sh
