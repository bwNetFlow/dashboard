FROM ubuntu:latest
RUN apt-get update ; apt-get install -y ca-certificates
WORKDIR /
ADD consumer /
ADD docker-init /
RUN chmod +x /consumer
ENV KAFKA_BROKERS="127.0.0.1:9092,[::1]:9092"
ENV KAFKA_TOPIC="flow-messages-enriched"
ENV KAFKA_CONSUMER_GROUP="dashboard"
ENV CONSUMER_FILTER_CUSTOMERID=""
ENV KAFKA_AUTH="false"
ENV KAFKA_TLS="false"
ENV KAFKA_USER=""
ENV KAFKA_PASS=""
ENV EXPORT_PROMETHEUS="true"
ENV EXPORT_INFLUX="false"
ENV EXPORT_INFLUX_URL=""
ENV EXPORT_INFLUX_USER=""
ENV EXPORT_INFLUX_PASS=""
ENV EXPORT_INFLUX_DATABASE=""
ENV EXPORT_INFLUX_FREQ=""
ENV EXPORT_INFLUX_PERCID=""
ENTRYPOINT [ "/docker-init" ]