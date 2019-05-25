package common

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"github.com/spf13/cobra"
)

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func TermSize() (uint, uint) {
	ws := &winsize{}
	retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		panic(errno)
	}
	return uint(ws.Col), uint(ws.Row)
}

func MetaNamespace(tenant string) string {
	if len(tenant) == 0 {
		return "default"
	}

	return fmt.Sprintf("cyclone-%s", tenant)
}

func ExecNamespace(tenant string) string {
	if len(tenant) == 0 {
		return "default"
	}

	return fmt.Sprintf("cyclone--%s", tenant)
}

func Equal(v1, v2 string) bool {
	return strings.TrimSpace(strings.ToLower(v1)) == strings.TrimSpace(strings.ToLower(v2))
}

func GetFlagValue(cmd *cobra.Command, name string) string {
	if f := cmd.Flag(name); f != nil {
		return f.Value.String()
	}

	return ""
}
