Zap the free space of a file system. Fill it with a secure byte pattern.

```
zapfsfs  -h  (or no arguments)
Show usage (this display).

zapfsfs  -t  TEMPDIR
Do a test run using the specified temp directory.

zapfsfs  [-b=N]  [-c=N]  [-n=N]  TEMPDIR
where
* TEMPDIR: Directory to create a temporary file (required)
* -b: Size of each buffer to write to a temp file (in MB); default: 10
* -c: Count of buffers to write in each pass (>0);  default: 1000
* -n: Number of file system freespace passes (>0); default: 5
```

