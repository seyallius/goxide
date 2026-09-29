// Command sysinfo prints the hardware/OS environment used for
// benchmark and profile documentation.
//
// Usage:
//
//	go run ./tools/sysinfo/ | tee docs/hardware.txt
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// ------------------------------------------- <Main> ------------------------------------------- //

func main() {
	fmt.Printf("os:     %s\n", runtime.GOOS)
	fmt.Printf("arch:   %s\n", runtime.GOARCH)
	fmt.Printf("go:     %s\n", runtime.Version())
	fmt.Printf("cpu:    %s\n", cpuModel())
	fmt.Printf("memory: %s\n", totalMemory())
	if h, err := os.Hostname(); err == nil {
		fmt.Printf("host:   %s\n", h)
	}
	if runtime.GOOS == "linux" {
		if k := kernelVersion(); k != "" {
			fmt.Printf("kernel: %s\n", k)
		}
	}
}

// -------------------------------------- Internal Helpers -------------------------------------- //

func cpuModel() string {
	switch runtime.GOOS {
	case "linux":
		if f, err := os.Open("/proc/cpuinfo"); err == nil {
			defer f.Close()
			sc := bufio.NewScanner(f)
			for sc.Scan() {
				line := sc.Text()
				if strings.HasPrefix(line, "model name") {
					if i := strings.Index(line, ":"); i >= 0 {
						return strings.TrimSpace(line[i+1:])
					}
				}
			}
		}
	case "darwin":
		return run("sysctl", "-n", "machdep.cpu.brand_string")
	case "windows":
		out := run("wmic", "cpu", "get", "name", "/value")
		for _, line := range strings.Split(out, "\n") {
			if strings.HasPrefix(line, "Name=") {
				return strings.TrimSpace(strings.TrimPrefix(line, "Name="))
			}
		}
	}
	return "unknown"
}

func totalMemory() string {
	switch runtime.GOOS {
	case "linux":
		if f, err := os.Open("/proc/meminfo"); err == nil {
			defer f.Close()
			sc := bufio.NewScanner(f)
			for sc.Scan() {
				line := sc.Text()
				if strings.HasPrefix(line, "MemTotal:") {
					fields := strings.Fields(line)
					if len(fields) >= 2 {
						if kb, err := strconv.ParseUint(fields[1], 10, 64); err == nil {
							return humanizeBytes(kb * 1024)
						}
					}
				}
			}
		}
	case "darwin":
		if out := run("sysctl", "-n", "hw.memsize"); out != "" {
			if b, err := strconv.ParseUint(strings.TrimSpace(out), 10, 64); err == nil {
				return humanizeBytes(b)
			}
		}
	case "windows":
		out := run("wmic", "computersystem", "get", "TotalPhysicalMemory", "/value")
		for _, line := range strings.Split(out, "\n") {
			if strings.HasPrefix(line, "TotalPhysicalMemory=") {
				v := strings.TrimSpace(strings.TrimPrefix(line, "TotalPhysicalMemory="))
				if b, err := strconv.ParseUint(v, 10, 64); err == nil {
					return humanizeBytes(b)
				}
			}
		}
	}
	return "unknown"
}

func kernelVersion() string {
	if b, err := os.ReadFile("/proc/version"); err == nil {
		fields := strings.Fields(string(b))
		if len(fields) >= 3 {
			return fields[2]
		}
	}
	return run("uname", "-r")
}

func run(name string, args ...string) string {
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func humanizeBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}
