FROM golang

WORKDIR /src
COPY . /src/
RUN make install

CMD sh scripts/start-all.sh
