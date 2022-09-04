# FSMon

FSMon is filesystem performance monitoring tool with Prometheus exporter API and Zabbix sender. All configuration
is read from enviroment variables.

## Install on Linux

```bash
mkdir /opt/fsmon
tar -C /opt/iwmon -xzf iwmon_v0.2.7_linux_amd64.tgz
cd /opt/fsmon
./iwmon -service install
./iwmon -service start
```

## Install on Windows

Unzip to C:\iwmon and run cmd.exe aith administrator account.

```cmd
cd c:\iwmon
iwmon.exe -service install
iwmon.exe -service start
```

## Enviroment variables

FS - path to filesystem that should be monitored
ZABBIX_SERVER - comma separated list of zabbix servers
PROMETHEUS - listening address + port to start prometheus exporter