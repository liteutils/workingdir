## WorkingDir

- 一款方便注册临时处理目录的小工具 可以设置定时清理策略

## How to Use

```

go get github.com/liteutils/workingdir

# init
WorkingDir := workingdir.New("")

# Create a temporary working directory located at /tmp/ofoov2 and continuously perform automatic cleanup of files older than 1 hour.
# 注册一个位于/tmp/2b的临时工作目录，并实时保持自动清理1小时之前的老文件
myDir,_ := WorkingDir.RegisterWorkingDir("ofoov2", 3600)


```