```go
alilog.InitFromConfigFile(os.Getenv("ALILOG_CONFIG"))
alilog.StartSlsLog()
defer alilog.CloseSlsLog(3000)
```
