# Development Time Series Databases & Dashboards

For developing the dashboard consumer, use the docker compose setup in this folder to deploy 

- Prometheus database, including the scraper which queries the dashboard consumer
- InfluxDB database where the dashboard consumer can push the data to
- Grafana and Chronograf dashboards to browse the data

Usage:

Within this folder, run `docker-compose up -d`. Then start the dashboard consumer.