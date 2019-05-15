FROM python:3.6-slim

WORKDIR /app

RUN apt-get update && apt-get install -y pandoc

ADD ./req.txt /app
RUN pip install -r ./req.txt
ADD . /app

CMD ["python", "server.py"]
