//+build !arm_logs

package main

import "syscall"

var dupfn = func(file int) { 
	syscall.Dup2(file, 2)
}
