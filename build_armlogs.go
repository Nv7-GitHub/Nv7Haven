//+build arm_logs

package main

import "syscall"

var dupfn = func(file int) {
	syscall.Dup3(file, 2, 0)
}
