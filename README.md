# Ping Monitor
A simple tool to monitor server pings written in Go and web technologies.

## How does it work?
It stores ping time in Redis storage every minute for each target (IP or server).
To increase the accuracy, it pings each target every 20 seconds and updates the result stored in the storage.
So even if one or two of the pings failed, there would be a ping time for that minute.

## Installation
```shell
git clone https://github.com/miladrahimi/ping-monitor.git
cd ping-monitor
cp .env.example .env
docker-compose up -d
docker-compose ps
```

## Configuration
Open `.env` with a text editor and change the environemnt variables.

Environemnt variables:
* **APP_EXPOSED_PORT**: The exposed port for web app
* **TARGETS**: The comma-separated list of targets (IP or server) to ping
* **TIMEZONE**: The timezone!

## Monitoring
Open your browser, surf localhost with the docker exposed port (default: 8585).

The chart is powered by [Chart.js](https://www.chartjs.org)

## Demo
<p align="center">
  <img alt="Demo" src="https://github.com/miladrahimi/ping-monitor/blob/master/demo.png?raw=true">
</p>

## See Also
* **[Health Monitor](https://github.com/miladrahimi/health-monitor)**: A simple tool to monitor server health

## License
Ping Monitor is initially created by [Milad Rahimi](https://miladrahimi.com)
and released under the [MIT License](http://opensource.org/licenses/mit-license.php).
