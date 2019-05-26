package common

import (
	"fmt"
	"math"
	"strconv"
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
	retCode, _, _ := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		return 0, 0
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

func PrefixMatch(origin, prefix string) bool {
	origin = strings.TrimSpace(strings.ToLower(origin))
	prefix = strings.TrimSpace(strings.ToLower(prefix))
	return strings.HasPrefix(origin, prefix)
}

func GetFlagValue(cmd *cobra.Command, name string) string {
	if f := cmd.Flag(name); f != nil {
		return f.Value.String()
	}

	return ""
}

func plural(count int, singular string) (result string) {
	result = strconv.Itoa(count) + singular
	return
}

func ReadableDuration(input int) (result string) {
	years := math.Floor(float64(input) / 60 / 60 / 24 / 7 / 30 / 12)
	seconds := input % (60 * 60 * 24 * 7 * 30 * 12)
	months := math.Floor(float64(seconds) / 60 / 60 / 24 / 7 / 30)
	seconds = input % (60 * 60 * 24 * 7 * 30)
	weeks := math.Floor(float64(seconds) / 60 / 60 / 24 / 7)
	seconds = input % (60 * 60 * 24 * 7)
	days := math.Floor(float64(seconds) / 60 / 60 / 24)
	seconds = input % (60 * 60 * 24)
	hours := math.Floor(float64(seconds) / 60 / 60)
	seconds = input % (60 * 60)
	minutes := math.Floor(float64(seconds) / 60)
	seconds = input % 60

	if years > 0 {
		result = ">1 year"
	} else if months > 0 {
		result = ">1 month"
	} else if weeks > 0 {
		result = plural(int(weeks), "w") + plural(int(days), "d") + plural(int(hours), "h") + plural(int(minutes), "m") + plural(int(seconds), "s")
	} else if days > 0 {
		result = plural(int(days), "d") + plural(int(hours), "h") + plural(int(minutes), "m") + plural(int(seconds), "s")
	} else if hours > 0 {
		result = plural(int(hours), "h") + plural(int(minutes), "m") + plural(int(seconds), "s")
	} else if minutes > 0 {
		result = plural(int(minutes), "m") + plural(int(seconds), "s")
	} else {
		result = plural(int(seconds), "s")
	}

	return
}
