FROM python:3.6-slim

WORKDIR /app

RUN apt-get update && apt-get install -y curl inotify-tools
RUN curl -sL https://github.com/jgm/pandoc/releases/download/2.7.2/pandoc-2.7.2-1-amd64.deb -o pandoc.deb  \
          && dpkg -i pandoc.deb

ADD ./req.txt /app
RUN pip install --no-cache -r ./req.txt
ADD ./server.py /app

CMD ["python", "server.py"]
