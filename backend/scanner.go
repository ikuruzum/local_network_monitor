package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type PortInfo struct {
	Port    int    `json:"port"`
	Process string `json:"process"`
}

func scanPorts() ([]PortInfo, error) {
	if runtime.GOOS == "windows" {
		return scanPortsWindows()
	}
	return scanPortsLinux()
}

func scanPortsLinux() ([]PortInfo, error) {
	var ports []PortInfo
	procNetTcp := "/proc/net/tcp"

	data, err := ioutil.ReadFile(procNetTcp)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")[1:]

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}

		portHex := strings.Split(fields[1], ":")[1]
		port, err := strconv.ParseInt(portHex, 16, 32)
		if err != nil {
			continue
		}

		inode := fields[9]
		process := findProcessByInodeLinux(inode)

		ports = append(ports, PortInfo{
			Port:    int(port),
			Process: process,
		})
	}

	return ports, nil
}

func scanPortsWindows() ([]PortInfo, error) {
	cmd := exec.Command("netstat", "-ano")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var ports []PortInfo
	lines := strings.Split(string(output), "\n")
	portMap := make(map[int]string)

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		if fields[0] != "TCP" {
			continue
		}

		localAddr := fields[1]
		parts := strings.Split(localAddr, ":")
		if len(parts) < 2 {
			continue
		}

		port, err := strconv.Atoi(parts[len(parts)-1])
		if err != nil {
			continue
		}

		pid := fields[len(fields)-1]
		if _, exists := portMap[port]; !exists {
			process := getProcessNameWindows(pid)
			portMap[port] = process
		}
	}

	for port, process := range portMap {
		ports = append(ports, PortInfo{
			Port:    port,
			Process: process,
		})
	}

	return ports, nil
}

func findProcessByInodeLinux(inode string) string {
	procDir := "/proc"
	files, _ := ioutil.ReadDir(procDir)

	for _, f := range files {
		if !f.IsDir() {
			continue
		}

		pid := f.Name()
		if _, err := strconv.Atoi(pid); err != nil {
			continue
		}

		fdPath := filepath.Join(procDir, pid, "fd")
		fds, err := ioutil.ReadDir(fdPath)
		if err != nil {
			continue
		}

		for _, fd := range fds {
			link, _ := os.Readlink(filepath.Join(fdPath, fd.Name()))
			if strings.Contains(link, "socket:["+inode+"]") {
				cmdline, _ := ioutil.ReadFile(filepath.Join(procDir, pid, "cmdline"))
				parts := strings.Split(string(cmdline), "\x00")
				if len(parts) > 0 {
					return filepath.Base(parts[0])
				}
			}
		}
	}

	return "unknown"
}

func getProcessNameWindows(pid string) string {
	cmd := exec.Command("tasklist", "/FI", "PID eq "+pid, "/FO", "CSV", "/NH")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}

	line := strings.TrimSpace(string(output))
	if line == "" {
		return "unknown"
	}

	parts := strings.Split(line, ",")
	if len(parts) > 0 {
		processName := strings.Trim(parts[0], "\"")
		return processName
	}

	return "unknown"
}