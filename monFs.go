package main

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/spetr/go-zabbix-sender"
)

func monFsUpdate(mountpoint, name string) {
	var (
		start           time.Time
		err             error
		testPath        string
		fh              *os.File
		buffer          []byte
		timeMonFsMkdir  float64
		timeMonFsList   float64
		timeMonFsCreate float64
		timeMonFsOpen   float64
		timeMonFsWrite  float64
		timeMonFsSync   float64
		timeMonFsRead   float64
		timeMonFsClose  float64
		timeMonFsStat   float64
		timeMonFsStatNx float64
		timeMonFsDelete float64
		timeMonFsRmdir  float64
	)

	for {
		func() {
			// Default values (no test / error in test)
			timeMonFsMkdir = -1
			timeMonFsList = -1
			timeMonFsCreate = -1
			timeMonFsOpen = -1
			timeMonFsWrite = -1
			timeMonFsSync = -1
			timeMonFsRead = -1
			timeMonFsClose = -1
			timeMonFsStat = -1
			timeMonFsStatNx = -1
			timeMonFsDelete = -1
			timeMonFsRmdir = -1

			if _, err = os.Stat(mountpoint); err != nil {
				logger.Errorf("Mountpoint %s path error: %s", name, err.Error())
				time.Sleep(10 * time.Second)
				return
			}

			testPath = path.Join(mountpoint, ".fsmon")
			// Create .fsmon folder and prepare data
			if err = os.MkdirAll(testPath, os.ModePerm); err != nil {
				logger.Errorf("Can not create %s fs testing directory: %s", name, err.Error())
				time.Sleep(10 * time.Second)
				return
			}
			testFolder := getRandString(16)
			testFile := fmt.Sprintf("%s.dat", getRandString(16))
			buffer = []byte(getRandString(8192))

			defer func() {
				if fh != nil {
					fh.Close()
				}
			}()

			// mkdir
			start = time.Now()
			if err = os.Mkdir(path.Join(testPath, testFolder), os.ModePerm); err != nil {
				return
			}
			timeMonFsMkdir = float64(time.Since(start).Microseconds())

			// list
			start = time.Now()
			if _, err = os.ReadDir(path.Join(testPath, testFolder)); err != nil {
				return
			}
			timeMonFsList = float64(time.Since(start).Microseconds())

			// create file
			start = time.Now()
			if fh, err = os.OpenFile(path.Join(testPath, testFolder, testFile), os.O_RDWR|os.O_CREATE, os.ModePerm); err != nil {
				return
			}
			timeMonFsCreate = float64(time.Since(start).Microseconds())
			fh.Close()

			// open file
			start = time.Now()
			if fh, err = os.OpenFile(path.Join(testPath, testFolder, testFile), os.O_RDWR, os.ModePerm); err != nil {
				return
			}
			timeMonFsOpen = float64(time.Since(start).Microseconds())

			// flock - TODO

			// write
			fh.SetWriteDeadline(time.Now().Add(2 * time.Second))
			start = time.Now()
			if _, err = fh.Write(buffer); err != nil {
				return
			}
			timeMonFsWrite = float64(time.Since(start).Microseconds())

			// sync
			fh.SetWriteDeadline(time.Now().Add(2 * time.Second))
			start = time.Now()
			if err = fh.Sync(); err != nil {
				return
			}
			timeMonFsSync = float64(time.Since(start).Microseconds())

			// read
			fh.SetReadDeadline(time.Now().Add(2 * time.Second))
			start = time.Now()
			if _, err = fh.ReadAt(buffer, 0); err != nil {
				return
			}
			timeMonFsRead = float64(time.Since(start).Microseconds())

			// close
			fh.SetWriteDeadline(time.Now().Add(2 * time.Second))
			start = time.Now()
			if err = fh.Close(); err != nil {
				return
			}
			timeMonFsClose = float64(time.Since(start).Microseconds())

			// stat
			start = time.Now()
			if _, err = os.Stat(path.Join(testPath, testFolder, testFile)); err != nil {
				return
			}
			timeMonFsStat = float64(time.Since(start).Microseconds())

			// statnx
			start = time.Now()
			_, _ = os.Stat(path.Join(testPath, testFolder, "non-existing.dat"))
			timeMonFsStatNx = float64(time.Since(start).Microseconds())

			// delete file
			start = time.Now()
			if err = os.Remove(path.Join(testPath, testFolder, testFile)); err != nil {
				return
			}
			timeMonFsDelete = float64(time.Since(start).Microseconds())

			// delete directory
			start = time.Now()
			if err = os.Remove(path.Join(testPath, testFolder)); err != nil {
				return
			}
			timeMonFsRmdir = float64(time.Since(start).Microseconds())

		}()

		// Zabbix Sender
		if len(conf.Zabbix.Servers) > 0 {
			var (
				metrics []*zabbix.Metric
				t       = time.Now().Unix()
			)
			metrics = append(metrics, zabbix.NewMetric(conf.Zabbix.Hostname, fmt.Sprintf("fsmon.mkdir.%s", name), fmt.Sprintf("%f", timeMonFsMkdir), true, t))
			metrics = append(metrics, zabbix.NewMetric(conf.Zabbix.Hostname, fmt.Sprintf("fsmon.list.%s", name), fmt.Sprintf("%f", timeMonFsList), true, t))
			metrics = append(metrics, zabbix.NewMetric(conf.Zabbix.Hostname, fmt.Sprintf("fsmon.create.%s", name), fmt.Sprintf("%f", timeMonFsCreate), true, t))
			metrics = append(metrics, zabbix.NewMetric(conf.Zabbix.Hostname, fmt.Sprintf("fsmon.open.%s", name), fmt.Sprintf("%f", timeMonFsOpen), true, t))
			metrics = append(metrics, zabbix.NewMetric(conf.Zabbix.Hostname, fmt.Sprintf("fsmon.write.%s", name), fmt.Sprintf("%f", timeMonFsWrite), true, t))
			metrics = append(metrics, zabbix.NewMetric(conf.Zabbix.Hostname, fmt.Sprintf("fsmon.sync.%s", name), fmt.Sprintf("%f", timeMonFsSync), true, t))
			metrics = append(metrics, zabbix.NewMetric(conf.Zabbix.Hostname, fmt.Sprintf("fsmon.read.%s", name), fmt.Sprintf("%f", timeMonFsRead), true, t))
			metrics = append(metrics, zabbix.NewMetric(conf.Zabbix.Hostname, fmt.Sprintf("fsmon.close.%s", name), fmt.Sprintf("%f", timeMonFsClose), true, t))
			metrics = append(metrics, zabbix.NewMetric(conf.Zabbix.Hostname, fmt.Sprintf("fsmon.stat.%s", name), fmt.Sprintf("%f", timeMonFsStat), true, t))
			metrics = append(metrics, zabbix.NewMetric(conf.Zabbix.Hostname, fmt.Sprintf("fsmon.statnx.%s", name), fmt.Sprintf("%f", timeMonFsStatNx), true, t))
			metrics = append(metrics, zabbix.NewMetric(conf.Zabbix.Hostname, fmt.Sprintf("fsmon.delete.%s", name), fmt.Sprintf("%f", timeMonFsDelete), true, t))
			metrics = append(metrics, zabbix.NewMetric(conf.Zabbix.Hostname, fmt.Sprintf("fsmon.rmdir.%s", name), fmt.Sprintf("%f", timeMonFsRmdir), true, t))
			for i := range conf.Zabbix.Servers {
				logger.Infof("Sending metrics to Zabbix server %s", conf.Zabbix.Servers[i].Host)
				zabbixSender := zabbix.NewSender(conf.Zabbix.Servers[i].Host)
				zabbixSender.ConnectTimeout = conf.Zabbix.Servers[i].ConnectTimeout
				zabbixSender.ReadTimeout = conf.Zabbix.Servers[i].ReadTimeout
				zabbixSender.WriteTimeout = conf.Zabbix.Servers[i].WriteTimeout
				if _, _, _, err := zabbixSender.SendMetrics(metrics); err != nil {
					logger.Errorf("Failed to send metrics to Zabbix server %s: %s", conf.Zabbix.Servers[i].Host, err)
				}
			}
		}

		time.Sleep(30 * time.Second)
	}
}
