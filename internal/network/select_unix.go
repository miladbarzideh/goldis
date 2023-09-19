//go:build !darwin

package network

import "syscall"

func platformSpecificSelect(n int, r *syscall.FdSet, w *syscall.FdSet, e *syscall.FdSet, timeout *syscall.Timeval) (err error) {
	_, err = syscall.Select(n, r, w, e, timeout)
	return err
}
