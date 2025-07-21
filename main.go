package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// Constants:
const mb uint64 = 1024 * 1024 // the Megabytes unit
const gb = 1024 * mb          // the Gigabytes unit

// Byte pattern used in buffer filling:
var pattern = []byte{0xAA, 0x55}

// Freespace quantity in bytes.
var tempFileSize uint64
var freeBytes uint64

// Command-line flags
var (
	bufferSizeMB = flag.Uint64("b", 10, "Size of each buffer to write to a temp file (in MB)")
	bufferCount  = flag.Uint64("c", 1000, "Count of buffers to write in each pass (>0)")
	npasses      = flag.Uint64("n", 5, "Number of file system freespace passes (>0)")
	testRun      = flag.Bool("t", false, "Number of file system freespace passes (>0)")
)

func usage() {
	fmt.Println()
	fmt.Println(`Zap the free space of a file system. Fill it with a secure byte pattern.

  Usage:

  zapfsfs  -h  (or no arguments)
    Show usage.

  zapfsfs  -t  TEMPDIR
	Do a test run using the specified temp directory.

  zapfsfs  [-b=N]  [-c=N]  [-n=N]  TEMPDIR
  where
    * TEMPDIR: Directory to create temporary file (required)
    * buffer_size: Size of each buffer to write to a temp file (in MB); default: 10
    * buffer_count: Count of buffers to write in each pass (>0);  default: 1000
    * npasses: Number of file system freespace passes (>0); default: 5
	`)
}

func bytesToGB(bytes uint64) float64 {
	return float64(bytes) / float64(gb)
}

func fillPattern(buf []byte, pattern []byte) {
	pLen := len(pattern)
	for ix := 0; ix < len(buf); ix += pLen {
		copy(buf[ix:], pattern)
	}
}

func scrubOnceSecure(tempDir string, bufferCount uint64, bufferSize uint64, pattern []byte) error {
	if len(pattern) < 1 {
		return fmt.Errorf("scrubOnceSecure: empty fill pattern array")
	}

	tempFile, err := os.CreateTemp(tempDir, "zapfsfs_*.tmp")
	if err != nil {
		return fmt.Errorf("scrubOnceSecure: create temp file: %w", err)
	}
	defer func() {
		_ = tempFile.Close()
		_ = os.Remove(tempFile.Name())
	}()

	buffer := make([]byte, bufferSize)
	fillPattern(buffer, pattern)

	for ix := uint64(0); ix < bufferCount; ix++ {
		if _, err := tempFile.Write(buffer); err != nil {
			return fmt.Errorf("scrubOnceSecure: write failed: %w", err)
		}
	}

	if err := tempFile.Sync(); err != nil {
		return fmt.Errorf("scrubOnceSecure: sync failed: %w", err)
	}

	return nil
}

func main() {

	var err error

	// Parse command line.
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		usage()
		return
	}

	// Get the temp directory path.
	if len(args) != 1 {
		_, _ = fmt.Fprintln(os.Stderr, "*** TEMPDIR is required as a positional argument.")
		usage()
		os.Exit(1)
	}
	tempDir := args[0]

	// Validation numeric parameters.
	if *bufferSizeMB < 1 || *bufferCount < 1 || *npasses < 1 {
		_, _ = fmt.Fprintln(os.Stderr, "*** buffer_size, buffer_count, and npasses must all be positive integers.")
		os.Exit(1)
	}
	bufferSize := *bufferSizeMB * mb
	tempFileSize = *bufferCount * bufferSize

	// Get freespace amount in bytes.
	freeBytes, err = getFreeSpace(tempDir)

	// Do a test run or a real run.
	if *testRun {
		fmt.Println("Begin, test run")
	} else {
		fmt.Printf("Begin, temp directory: %s, Buffer count: %d, Buffer size: %d MB, temp file size: %.2f GB, #passes: %d\n",
			tempDir, *bufferCount, *bufferSizeMB, bytesToGB(tempFileSize), *npasses)
		fmt.Printf("Filesystem free space: %.2f GB\n", bytesToGB(freeBytes))
		// Make sure that the filesystem has enough freespace.
		if freeBytes < tempFileSize {
			_, _ = fmt.Fprintf(os.Stderr, "*** Not enough free space in temp directory file system (%.2f GB).\n",
				bytesToGB(freeBytes))
			_, _ = fmt.Fprintln(os.Stderr, "*** Modify one or more numeric parameters and try again.")
			os.Exit(1)
		}
	}

	// Ensure the temp directory exists.
	stat, err := os.Stat(tempDir)
	if err != nil || !stat.IsDir() {
		_, _ = fmt.Fprintf(os.Stderr, "'%s' must be an existing directory that is writable.\n", tempDir)
		os.Exit(1)
	}

	// Ensure the temp directory is writable.
	dummyPath := filepath.Join(tempDir, "zapfsfs_dummy.bin")
	dummyFile, err := os.Create(dummyPath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "os.Create(%s) failed, err: %v.\n", dummyPath, err)
		os.Exit(1)
	}
	_, err = dummyFile.Write(pattern)
	_ = dummyFile.Close()
	_ = os.Remove(dummyPath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "dummy.Write(%s) failed, err: %v.\n", dummyPath, err)
		os.Exit(1)
	}

	// If not a test run, do npasses of secure free space scrubbing.
	if !*testRun {
		for ix := uint64(0); ix < *npasses; ix++ {
			pass := ix + 1
			fmt.Printf("Pass %d of %d...\n", pass, *npasses)
			err = scrubOnceSecure(tempDir, *bufferCount, bufferSize, pattern)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error during pass %d, err: %v\n", pass, err)
				os.Exit(1)
			}
		}
	}

	// Success!
	fmt.Println("End.")
}
