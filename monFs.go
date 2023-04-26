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
		timeMonFsMkdir  int64
		timeMonFsList   int64
		timeMonFsCreate int64
		timeMonFsOpen   int64
		timeMonFsWrite  int64
		timeMonFsSync   int64
		timeMonFsRead   int64
		timeMonFsClose  int64
		timeMonFsStat   int64
		timeMonFsStatNx int64
		timeMonFsDelete int64
		timeMonFsRmdir  int64
	)

	// Run zabbix discovery
	metrics := make([]*zabbix.Metric, 1)
	t := time.Now().Unix()
	metrics[0] = zabbix.NewMetric(conf.Zabbix.Hostname, "fsmon.discovery", fmt.Sprintf("[{\"{#NAME}\": \"%s\"}]", name), true, t)
	for i := range conf.Zabbix.Servers {
		logger.Infof("Starting service discovery for host %s on %s with filesystem %s", conf.Zabbix.Hostname, conf.Zabbix.Servers[i].Host, name)
		zabbixSender := zabbix.NewSender(conf.Zabbix.Servers[i].Host)
		zabbixSender.ConnectTimeout = conf.Zabbix.Servers[i].ConnectTimeout
		zabbixSender.ReadTimeout = conf.Zabbix.Servers[i].ReadTimeout
		zabbixSender.WriteTimeout = conf.Zabbix.Servers[i].WriteTimeout
		zabbixResponse, err, _, _ := zabbixSender.SendMetrics(metrics)
		if *debugFlag {
			logger.Infof("Zabbix response info: %s", zabbixResponse.Info)
			logger.Infof("Zabbix response: %s", zabbixResponse.Response)
		}
		if err != nil {
			logger.Errorf("Failed to send metxrics to %s: %s", conf.Zabbix.Servers[i].Host, err)
		}
	}

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
			timeMonFsMkdir = time.Since(start).Microseconds()
			if *debugFlag {
				logger.Infof("Measured time mkdir on %s: %dμs", name, timeMonFsMkdir)
			}

			// list
			start = time.Now()
			if _, err = os.ReadDir(path.Join(testPath, testFolder)); err != nil {
				return
			}
			timeMonFsList = time.Since(start).Microseconds()
			if *debugFlag {
				logger.Infof("Measured time list on %s: %dμs", name, timeMonFsList)
			}

			// create file
			start = time.Now()
			if fh, err = os.OpenFile(path.Join(testPath, testFolder, testFile), os.O_RDWR|os.O_CREATE, os.ModePerm); err != nil {
				return
			}
			timeMonFsCreate = time.Since(start).Microseconds()
			fh.Close()
			if *debugFlag {
				logger.Infof("Measured time create on %s: %dμs", name, timeMonFsCreate)
			}

			// open file
			start = time.Now()
			if fh, err = os.OpenFile(path.Join(testPath, testFolder, testFile), os.O_RDWR, os.ModePerm); err != nil {
				return
			}
			timeMonFsOpen = time.Since(start).Microseconds()
			if *debugFlag {
				logger.Infof("Measured time open on %s: %dμs", name, timeMonFsOpen)
			}

			// flock - TODO

			// write
			fh.SetWriteDeadline(time.Now().Add(2 * time.Second))
			start = time.Now()
			if _, err = fh.Write(buffer); err != nil {
				return
			}
			timeMonFsWrite = time.Since(start).Microseconds()
			if *debugFlag {
				logger.Infof("Measured time write on %s: %dμs", name, timeMonFsWrite)
			}

			// sync
			fh.SetWriteDeadline(time.Now().Add(2 * time.Second))
			start = time.Now()
			if err = fh.Sync(); err != nil {
				return
			}
			timeMonFsSync = time.Since(start).Microseconds()
			if *debugFlag {
				logger.Infof("Measured time sync on %s: %dμs", name, timeMonFsSync)
			}

			// read
			fh.SetReadDeadline(time.Now().Add(2 * time.Second))
			start = time.Now()
			if _, err = fh.ReadAt(buffer, 0); err != nil {
				return
			}
			timeMonFsRead = time.Since(start).Microseconds()
			if *debugFlag {
				logger.Infof("Measured time read on %s: %dμs", name, timeMonFsRead)
			}

			// close
			fh.SetWriteDeadline(time.Now().Add(2 * time.Second))
			start = time.Now()
			if err = fh.Close(); err != nil {
				return
			}
			timeMonFsClose = time.Since(start).Microseconds()
			if *debugFlag {
				logger.Infof("Measured time close on %s: %dμs", name, timeMonFsClose)
			}

			// stat
			start = time.Now()
			if _, err = os.Stat(path.Join(testPath, testFolder, testFile)); err != nil {
				return
			}
			timeMonFsStat = time.Since(start).Microseconds()
			if *debugFlag {
				logger.Infof("Measured time stat on %s: %dμs", name, timeMonFsStat)
			}

			// statnx
			start = time.Now()
			_, _ = os.Stat(path.Join(testPath, testFolder, "non-existing.dat"))
			timeMonFsStatNx = time.Since(start).Microseconds()
			if *debugFlag {
				logger.Infof("Measured time statnx on %s: %dμs", name, timeMonFsStatNx)
			}

			// delete file
			start = time.Now()
			if err = os.Remove(path.Join(testPath, testFolder, testFile)); err != nil {
				return
			}
			timeMonFsDelete = time.Since(start).Microseconds()
			if *debugFlag {
				logger.Infof("Measured time delete on %s: %dμs", name, timeMonFsDelete)
			}

			// delete directory
			start = time.Now()
			if err = os.Remove(path.Join(testPath, testFolder)); err != nil {
				return
			}
			timeMonFsRmdir = time.Since(start).Microseconds()
			if *debugFlag {
				logger.Infof("Measured time rmdir on %s: %dμs", name, timeMonFsRmdir)
			}

		}()

		// Zabbix Sender
		if len(conf.Zabbix.Servers) > 0 {
			metrics := make([]*zabbix.Metric, 12)
			t := time.Now().Unix()
			metrics[0] = zabbix.NewMetric(conf.Zabbix.Hostname, fmt.Sprintf("fsmon.mkdir.[%s]", name), fmt.Sprintf("%d", timeMonFsMkdir), true, t)
			metrics[1] = zabbix.NewMetric(conf.Zabbix.Hostname, fmt.Sprintf("fsmon.list.[%s]", name), fmt.Sprintf("%d", timeMonFsList), true, t)
			metrics[2] = zabbix.NewMetric(conf.Zabbix.Hostname, fmt.Sprintf("fsmon.create.[%s]", name), fmt.Sprintf("%d", timeMonFsCreate), true, t)
			metrics[3] = zabbix.NewMetric(conf.Zabbix.Hostname, fmt.Sprintf("fsmon.open.[%s]", name), fmt.Sprintf("%d", timeMonFsOpen), true, t)
			metrics[4] = zabbix.NewMetric(conf.Zabbix.Hostname, fmt.Sprintf("fsmon.write.[%s]", name), fmt.Sprintf("%d", timeMonFsWrite), true, t)
			metrics[5] = zabbix.NewMetric(conf.Zabbix.Hostname, fmt.Sprintf("fsmon.sync.[%s]", name), fmt.Sprintf("%d", timeMonFsSync), true, t)
			metrics[6] = zabbix.NewMetric(conf.Zabbix.Hostname, fmt.Sprintf("fsmon.read.[%s]", name), fmt.Sprintf("%d", timeMonFsRead), true, t)
			metrics[7] = zabbix.NewMetric(conf.Zabbix.Hostname, fmt.Sprintf("fsmon.close.[%s]", name), fmt.Sprintf("%d", timeMonFsClose), true, t)
			metrics[8] = zabbix.NewMetric(conf.Zabbix.Hostname, fmt.Sprintf("fsmon.stat.[%s]", name), fmt.Sprintf("%d", timeMonFsStat), true, t)
			metrics[9] = zabbix.NewMetric(conf.Zabbix.Hostname, fmt.Sprintf("fsmon.statnx.[%s]", name), fmt.Sprintf("%d", timeMonFsStatNx), true, t)
			metrics[10] = zabbix.NewMetric(conf.Zabbix.Hostname, fmt.Sprintf("fsmon.delete.[%s]", name), fmt.Sprintf("%d", timeMonFsDelete), true, t)
			metrics[11] = zabbix.NewMetric(conf.Zabbix.Hostname, fmt.Sprintf("fsmon.rmdir.[%s]", name), fmt.Sprintf("%d", timeMonFsRmdir), true, t)
			for i := range conf.Zabbix.Servers {
				logger.Infof("Sending metrics to Zabbix server %s", conf.Zabbix.Servers[i].Host)
				zabbixSender := zabbix.NewSender(conf.Zabbix.Servers[i].Host)
				zabbixSender.ConnectTimeout = conf.Zabbix.Servers[i].ConnectTimeout
				zabbixSender.ReadTimeout = conf.Zabbix.Servers[i].ReadTimeout
				zabbixSender.WriteTimeout = conf.Zabbix.Servers[i].WriteTimeout
				zabbixResponse, err, _, _ := zabbixSender.SendMetrics(metrics)
				if *debugFlag {
					logger.Infof("Zabbix response info: %s", zabbixResponse.Info)
					logger.Infof("Zabbix response: %s", zabbixResponse.Response)
				}
				if err != nil {
					logger.Errorf("Failed to send metrics to Zabbix server %s: %s", conf.Zabbix.Servers[i].Host, err)
				}
			}
		}

		time.Sleep(30 * time.Second)
	}
}
