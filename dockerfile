FROM ubuntu:latest
WORKDIR /
ADD processor /
ADD docker-init /
RUN chmod +x /processor
ENV KAFKA_BROKERS="127.0.0.1:9092,[::1]:9092"
ENV KAFKA_TOPIC="flow-messages-enriched"
ENV KAFKA_CONSUMER_GROUP="dashboard"
ENV CONSUMER_FILTER_CUSTOMERID="10109"
ENTRYPOINT [ "/docker-init" ]