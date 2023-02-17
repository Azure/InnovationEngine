FROM ubuntu:20.04
ENV TERM xterm
WORKDIR $HOME

RUN apt-get update && \
    apt-get install -y python3 python3-pip git

WORKDIR /app

COPY . /app
RUN chmod +x main.py
RUN pip3 install -r requirements.txt
CMD [ "python3", "main.py"]
