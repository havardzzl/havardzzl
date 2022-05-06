package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"math"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"syscall"
	"time"

	"github.com/opencontainers/runc/libcontainer/cgroups"
	"github.com/opencontainers/runc/libcontainer/cgroups/fscommon"
	"k8s.io/klog/v2"
)

const (
	baseDir      = "/sys/fs/cgroup/cpu/kubepods.slice"
	quotaFile    = "cpu.cfs_quota_us"
	burstFile    = "cpu.cfs_burst_us"
	defaultQuota = 800000
)

var (
	changeFiles int

	// pod: kubepods-burstable-pod{$pod-uid}.slice
	// container: docker-{$ContainerID}.scope
	podRegexp       = regexp.MustCompile(`kubepods-burstable-pod([^\}]+).slice`)
	containerRegexp = regexp.MustCompile(`docker-([^\}]+).scope`)

	// container维度的
	metricsKV map[metricKey]metricValue
)

func getCpuStats(path string, stats *metricValue) error {
	const file = "cpu.stat"
	f, err := cgroups.OpenFile(path, file, os.O_RDONLY)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		t, v, err := fscommon.ParseKeyValue(sc.Text())
		if err != nil {
			return &fscommon.ParseError{Path: path, File: file, Err: err}
		}
		switch t {
		case "nr_periods":
			stats.periods = v
		case "nr_throttled":
			stats.throttledPeriods = v
		case "throttled_time":
			stats.throttledTime = v
		case "nr_bursts":
			stats.burstPeriods = v
		case "burst_time":
			stats.burstTime = v
		}
	}
	return nil
}

func walkDirFunc(path string, d fs.DirEntry, err error) error {
	if err != nil || d == nil {
		klog.ErrorS(err, "walkDirFunc failed: ", path)
		return nil
	}
	if !d.IsDir() {
		return nil
	}
	burstFileName := filepath.Join(path, burstFile)
	_, errl := os.Lstat(burstFileName)
	if errl != nil {
		return nil
	}
	quota, errq := fscommon.GetCgroupParamInt(path, quotaFile)
	if errq != nil {
		if !errors.Is(errq, os.ErrNotExist) {
			klog.ErrorS(errq, "GetCgroupParamInt failed: ", path, quotaFile)
		}
		return nil
	}
	if quota < 0 || quota == math.MaxInt64 {
		cgroups.WriteFile(path, burstFile, fmt.Sprintf("%d", defaultQuota))
		return nil
	}
	// burst的值是quota的2倍效果更好，但是测试发现不行
	errb := cgroups.WriteFile(path, burstFile, fmt.Sprintf("%d", quota))
	if errb != nil {
		klog.ErrorS(errb, "WriteFile failed: ", path, burstFileName, quota)
	} else {
		changeFiles++
	}
	// 上报metrics
	pod := filepath.Base(path)
	container := d.Name()
	podMatch := podRegexp.FindStringSubmatch(pod)
	containerMatch := containerRegexp.FindStringSubmatch(container)
	if len(podMatch) == 2 && len(containerMatch) == 2 {
		key := metricKey{podMatch[1], containerMatch[1]}
		value := metricValue{}
		if err := getCpuStats(path, &value); err != nil {
			klog.ErrorS(err, "GetStats failed: ", path)
		}
		metricsKV[key] = value
	}

	return nil
}

func work() {
	defer func() {
		if err := recover(); err != nil {
			klog.Error("work panic: ", err)
		}
	}()
	for k := range metricsKV {
		delete(metricsKV, k)
	}
	changeFiles = 0
	filepath.WalkDir(baseDir, walkDirFunc)
	klog.V(5).Infof("metricsKV: %+v", metricsKV)
	for k, v := range metricsKV {
		RecordMetrics(k, v)
	}
	if changeFiles == 0 {
		klog.Info("没有发现需要修改的文件")
		return
	}
	klog.Info("changeFiles: ", changeFiles)
}

func main() {
	go func() {
		for {
			work()
			time.Sleep(5 * time.Second)
		}
	}()
	StartMetricsServer()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-c
	StopMetricsServer()
	fmt.Println("byebye")
	os.Exit(0)
}
