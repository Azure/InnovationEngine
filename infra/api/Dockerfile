FROM mcr.microsoft.com/cbl-mariner/base/core:2.0

ARG HOST=0.0.0.0
ARG PORT=8080

WORKDIR /api

RUN tdnf update && \
  tdnf install golang make ca-certificates -y

COPY . .

RUN make build-api

EXPOSE 8080

CMD ["./bin/api"]