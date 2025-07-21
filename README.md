Zap the free space of a file system. Fill it with a secure byte pattern.

```
zapfsfs  -h  (or no arguments)
Show usage (this display).

zapfsfs  -t  TEMPDIR
Do a test run using the specified temp directory.

zapfsfs  [-b=N]  [-c=N]  [-n=N]  TEMPDIR
where
* TEMPDIR: Directory to create a temporary file (required)
* buffer_size: Size of each buffer to write to a temp file (in MB); default: 10
* buffer_count: Count of buffers to write in each pass (>0);  default: 1000
* npasses: Number of file system freespace passes (>0); default: 5
```

