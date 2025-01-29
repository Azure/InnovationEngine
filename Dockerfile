FROM golang:1.20-alpine

RUN apk update
RUN apk add --no-cache --update \
    bash \
    cargo \
    git \
    gcc \
    libffi-dev \
    make \
    musl-dev \
    openssl-dev \
    python3 \
    py3-pip \
    python3-dev

WORKDIR /InnovationEngine

# Create a virtual environment and install the experimental Authoring Tools and az cli
RUN python3 -m venv /InnovationEngine/venv
RUN /InnovationEngine/venv/bin/pip install openai azure-identity requests pygithub
RUN /InnovationEngine/venv/bin/pip install azure-cli

ENV VIRTUAL_ENV=/InnovationEngine/venv
ENV PATH="$VIRTUAL_ENV/bin:$PATH"

RUN mkdir -p AuthoringTools
RUN wget -O AuthoringTools/ada.py https://raw.githubusercontent.com/naman-msft/exec/main/tools/ada.py
RUN chmod +x AuthoringTools/ada.py

# Install the Innovation Engine
COPY . .
RUN make build-ie
ENV PATH="/InnovationEngine/bin:${PATH}"

CMD ["sh", "-c", "ie execute docs/helloWorldDemo.md"]