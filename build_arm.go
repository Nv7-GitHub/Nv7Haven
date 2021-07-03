//+build arm

package main

import "syscall"

var dupfn = syscall.Dup3
