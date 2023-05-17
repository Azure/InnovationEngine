FROM mcr.microsoft.com/cbl-mariner/base/core:2.0

ARG HOST=0.0.0.0
ARG PORT=8080

WORKDIR /api

RUN tdnf install golang make -y

COPY . .

RUN make build-runner

CMD ["./bin/runner"]