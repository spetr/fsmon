# FSMon

FSMon is filesystem performance monitoring tool with Prometheus exporter API and Zabbix sender. All configuration
is read from enviroment variables.

## Install on Linux

```bash
mkdir /opt/fsmon
tar -C /opt/fsmon -xzf fsmon_v0.2.7_linux_amd64.tgz
cd /opt/fsmon
./fsmon -service install
./fsmon -service start
```

## Install on Windows

Unzip to C:\fsmon and run cmd.exe aith administrator account.

```cmd
cd c:\fsmon
fsmon.exe -service install
fsmon.exe -service start
```
