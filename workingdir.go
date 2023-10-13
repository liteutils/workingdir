package workingdir

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"time"
)

type WorkingDir struct {
	BaseDir  string
	strategy []*ClearStrategy
}
type ClearStrategy struct {
	dirPath       string
	purgeInterval int64      // 清理间隔 单位秒
	purgeTime     *time.Time // 上次清理时间
}

func New(baseDir string) *WorkingDir {
	if baseDir == "" {
		baseDir = os.TempDir()
	}
	wd := &WorkingDir{
		BaseDir: baseDir,
	}
	go wd.PurgeWorkingDir()
	return wd
}

func (w *WorkingDir) fileExist(path string) bool {
	fi, err := os.Lstat(path)
	if err == nil {
		return !fi.IsDir()
	}
	return !os.IsNotExist(err)
}

// 注册新的工作目录
func (w *WorkingDir) RegisterWorkingDir(dirName string, purgeIntervalSeconds int64) (string, error) {
	var workingDir string = path.Join(w.BaseDir, dirName)
	if !w.fileExist(workingDir) {
		if err := os.MkdirAll(workingDir, os.ModePerm);err != nil {
			return "", err
		}
	}
	w.strategy = append(w.strategy, &ClearStrategy{
		dirPath:       workingDir,
		purgeInterval: purgeIntervalSeconds,
		purgeTime:     nil,
	})
	return workingDir, nil
}

// 定时清理已过期缓存文件
func (w *WorkingDir) PurgeWorkingDir() {
	if runtime.GOOS == "windows" {
		return
	}
	for {
		for _, strategy := range w.strategy {
			if strategy.purgeTime == nil || time.Now().Unix()-strategy.purgeTime.Unix() >= strategy.purgeInterval {
				cmd := exec.Command("bash", "-c", fmt.Sprintf("find "+path.Join(strategy.dirPath, "*")+" -type f,d -mmin +%d -exec rm {} \\;", int(strategy.purgeInterval/60)))
				cmd.Run()
				strategy.purgeTime = w.GetTime(time.Now())
			}
		}
		time.Sleep(time.Minute)
	}
}
func (w *WorkingDir) GetTime(t time.Time) *time.Time {
	return &t
}
