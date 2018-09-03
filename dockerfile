FROM alpine:latest
WORKDIR /
ADD processor /
RUN chmod +x /processor
ENV KAFKA_BROKERS="127.0.0.1:9092,[::1]:9092"
ENV KAFKA_TOPIC="flow-messages-enriched"
ENV KAFKA_CONSUMER_GROUP="dashboard"
CMD /processor --kafka.brokers \"$KAFKA_BROKERS\" \
  --kafka.topic \"$KAFKA_TOPIC\" \
  --kafka.consumer_group \"$KAFKA_CONSUMER_GROUP\"