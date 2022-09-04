package main

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/spetr/go-zabbix-sender"
)

var (
	monFsMkdir = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "fs_mail_mkdir",
			Help: "Filesystem mail storage - mkdir (microseconds)",
		},
	)
	monFsList = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "fs_mail_list",
			Help: "Filesystem mail storage - list (microseconds)",
		},
	)
	monFsCreate = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "fs_mail_create",
			Help: "Filesystem mail storage - create (microseconds)",
		},
	)
	monFsOpen = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "fs_mail_open",
			Help: "Filesystem mail storage - open (microseconds)",
		},
	)
	monFsWrite = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "fs_mail_write",
			Help: "Filesystem mail storage - write (microseconds)",
		},
	)
	monFsSync = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "fs_mail_sync",
			Help: "Filesystem mail storage - sync (microseconds)",
		},
	)
	monFsRead = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "fs_mail_read",
			Help: "Filesystem mail storage - read (microseconds)",
		},
	)
	monFsClose = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "fs_mail_close",
			Help: "Filesystem mail storage - close (microseconds)",
		},
	)
	monFsStat = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "fs_mail_stat",
			Help: "Filesystem mail storage - stat (microseconds)",
		},
	)
	monFsStatNx = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "fs_mail_statnx",
			Help: "Filesystem mail storage - statnx (microseconds)",
		},
	)
	monFsDelete = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "fs_mail_delete",
			Help: "Filesystem mail storage - delete (microseconds)",
		},
	)
	monFsRmdir = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "fs_mail_rmdir",
			Help: "Filesystem mail storage - rmdir (microseconds)",
		},
	)
)

func monFsUpdate(r *prometheus.Registry) {

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

	if conf.Prometheus != "" {
		r.MustRegister(monFsMkdir)
		r.MustRegister(monFsList)
		r.MustRegister(monFsCreate)
		r.MustRegister(monFsOpen)
		r.MustRegister(monFsWrite)
		r.MustRegister(monFsSync)
		r.MustRegister(monFsRead)
		r.MustRegister(monFsClose)
		r.MustRegister(monFsStat)
		r.MustRegister(monFsStatNx)
		r.MustRegister(monFsDelete)
		r.MustRegister(monFsRmdir)
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

			if _, err = os.Stat(conf.Filesystem); err != nil {
				logger.Errorf("Mail path error: %s", err.Error())
				time.Sleep(10 * time.Second)
				return
			}

			testPath = path.Join(conf.Filesystem, ".fsmon")
			// Create .fsmon folder and prepare data
			if err = os.MkdirAll(testPath, os.ModePerm); err != nil {
				logger.Errorf("Can not create mail fs testing directort: %s", err.Error())
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

		// Prometheus Exporter
		if conf.Prometheus != "" {
			monFsMkdir.Set(timeMonFsMkdir)
			monFsList.Set(timeMonFsList)
			monFsCreate.Set(timeMonFsCreate)
			monFsOpen.Set(timeMonFsOpen)
			monFsWrite.Set(timeMonFsWrite)
			monFsSync.Set(timeMonFsSync)
			monFsRead.Set(timeMonFsRead)
			monFsClose.Set(timeMonFsClose)
			monFsStat.Set(timeMonFsStat)
			monFsStatNx.Set(timeMonFsStatNx)
			monFsDelete.Set(timeMonFsDelete)
			monFsRmdir.Set(timeMonFsRmdir)
		}

		// Zabbix Sender
		if len(conf.ZabbixServer) > 0 {
			var (
				metrics     []*zabbix.Metric
				t           = time.Now().Unix()
				hostname, _ = os.Hostname()
			)
			metrics = append(metrics, zabbix.NewMetric(hostname, "fs.mail_mkdir", fmt.Sprintf("%f", timeMonFsMkdir), true, t))
			metrics = append(metrics, zabbix.NewMetric(hostname, "fs.mail_list", fmt.Sprintf("%f", timeMonFsList), true, t))
			metrics = append(metrics, zabbix.NewMetric(hostname, "fs.mail_create", fmt.Sprintf("%f", timeMonFsCreate), true, t))
			metrics = append(metrics, zabbix.NewMetric(hostname, "fs.mail_open", fmt.Sprintf("%f", timeMonFsOpen), true, t))
			metrics = append(metrics, zabbix.NewMetric(hostname, "fs.mail_write", fmt.Sprintf("%f", timeMonFsWrite), true, t))
			metrics = append(metrics, zabbix.NewMetric(hostname, "fs.mail_sync", fmt.Sprintf("%f", timeMonFsSync), true, t))
			metrics = append(metrics, zabbix.NewMetric(hostname, "fs.mail_read", fmt.Sprintf("%f", timeMonFsRead), true, t))
			metrics = append(metrics, zabbix.NewMetric(hostname, "fs.mail_close", fmt.Sprintf("%f", timeMonFsClose), true, t))
			metrics = append(metrics, zabbix.NewMetric(hostname, "fs.mail_stat", fmt.Sprintf("%f", timeMonFsStat), true, t))
			metrics = append(metrics, zabbix.NewMetric(hostname, "fs.mail_statnx", fmt.Sprintf("%f", timeMonFsStatNx), true, t))
			metrics = append(metrics, zabbix.NewMetric(hostname, "fs.mail_delete", fmt.Sprintf("%f", timeMonFsDelete), true, t))
			metrics = append(metrics, zabbix.NewMetric(hostname, "fs.mail_rmdir", fmt.Sprintf("%f", timeMonFsRmdir), true, t))
			for i := range conf.ZabbixServer {
				zabbix.NewSender(conf.ZabbixServer[i]).SendMetrics(metrics)
			}
		}

		time.Sleep(30 * time.Second)
	}
}
