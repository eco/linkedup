FROM golang

RUN apt-get update && apt-get install -y awscli

WORKDIR /src

# Cache dependencies to speed up rebuilds
COPY go.mod go.sum ./
RUN go mod download

COPY . /src/
RUN make clean all

CMD sh scripts/start-all.sh
