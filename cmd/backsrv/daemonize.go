package main

import (
	"fmt"
	"os"
	"syscall"
)

func daemon(nochdir, noclose int) (int, error) {
	// already a daemon
	if syscall.Getppid() == 1 {
		/* Change the file mode mask */
		syscall.Umask(0)

		if nochdir == 0 {
			err := os.Chdir("/")
			if err != nil {
				return 0, err
			}
		}

		return 0, nil
	}

	files := make([]*os.File, 3, 6)
	if noclose == 0 {
		nullDev, err := os.OpenFile("/dev/null", 0, 0)
		if err != nil {
			return 1, err
		}
		files[0], files[1], files[2] = nullDev, nullDev, nullDev
	} else {
		files[0], files[1], files[2] = os.Stdin, os.Stdout, os.Stderr
	}

	dir, _ := os.Getwd()
	sysAttrs := syscall.SysProcAttr{Setsid: true}
	attrs := os.ProcAttr{Dir: dir, Env: os.Environ(), Files: files, Sys: &sysAttrs}

	proc, err := os.StartProcess(os.Args[0], os.Args, &attrs)
	if err != nil {
		return -1, fmt.Errorf("can't create process %s: %s", os.Args[0], err)
	}
	if err := proc.Release(); err != nil {
		return 0, err
	}
	os.Exit(0)

	return 0, nil
}
