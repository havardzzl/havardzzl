package main

import (
	"errors"
	"fmt"
	"io/fs"
	"math"
	"os"
	"path/filepath"

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
)

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
		cgroups.WriteFile(path, burstFileName, fmt.Sprintf("%d", defaultQuota))
		return nil
	}
	// burst的值是quota的2倍效果更好，但是测试发现不行
	errb := cgroups.WriteFile(path, burstFileName, fmt.Sprintf("%d", quota))
	if errb != nil {
		klog.ErrorS(errb, "WriteFile failed: ", path, burstFileName, quota)
	} else {
		changeFiles++
	}
	return nil
}

func work() {
	defer func() {
		if err := recover(); err != nil {
			klog.Error("work panic: ", err)
		}
	}()
	changeFiles = 0
	filepath.WalkDir(baseDir, walkDirFunc)
	if changeFiles == 0 {
		klog.Info("没有发现需要修改的文件")
		return
	}
	klog.Info("changeFiles: ", changeFiles)
}
