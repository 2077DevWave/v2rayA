package v2ray

import (
	"V2RayA/global"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func EnableV2rayService() (err error) {
	switch global.ServiceControlMode {
	case global.DockerMode, global.CommonMode: //docker, common中无需enable service
	case global.ServiceMode:
		_, err = exec.Command("sh", "-c", "update-rc.d v2ray enable").CombinedOutput()
	case global.SystemctlMode:
		_, err = exec.Command("sh", "-c", "systemctl enable v2ray").Output()
	}
	return
}

func DisableV2rayService() (err error) {
	switch global.ServiceControlMode {
	case global.DockerMode, global.CommonMode: //docker, common中无需disable service
	case global.ServiceMode:
		_, err = exec.Command("sh", "-c", "update-rc.d v2ray disable").CombinedOutput()
	case global.SystemctlMode:
		_, err = exec.Command("sh", "-c", "systemctl disable v2ray").Output()
	}
	return
}

func GetV2rayServiceFilePath() (path string, err error) {
	var out []byte

	if global.ServiceControlMode == global.SystemctlMode {
		out, err = exec.Command("sh", "-c", "systemctl status v2ray|grep Loaded|awk '{print $3}'").Output()
		if err != nil {
			path = `/usr/lib/systemd/system/v2ray.service`
		}
	} else if global.ServiceControlMode == global.ServiceMode {
		out, err = exec.Command("sh", "-c", "service v2ray status|grep Loaded|awk '{print $3}'").Output()
		if err != nil || strings.TrimSpace(string(out)) == "(Reason:" {
			path = `/lib/systemd/system/v2ray.service`
		}
	} else {
		err = errors.New("当前环境无法使用systemctl和service命令")
		return
	}
	sout := strings.TrimSpace(string(out))
	path = sout[1 : len(sout)-1]
	return
}

func LiberalizeProcFile() (err error) {
	if global.ServiceControlMode != global.SystemctlMode && global.ServiceControlMode != global.ServiceMode {
		return
	}
	p, err := GetV2rayServiceFilePath()
	if err != nil {
		return
	}
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return
	}
	s := string(b)
	if strings.Contains(s, "LimitNPROC=500") && strings.Contains(s, "LimitNOFILE=1000000") {
		return
	}
	lines := strings.Split(s, "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.HasPrefix(lines[i], "LimitNPROC=") || strings.HasPrefix(lines[i], "LimitNOFILE=") {
			lines = append(lines[:i], lines[i+1:]...)
		}
	}
	for i, line := range lines {
		if strings.ToLower(line) == "[service]" {
			s = strings.Join(lines[:i+1], "\n")
			s += "\nLimitNPROC=500\nLimitNOFILE=1000000\n"
			s += strings.Join(lines[i+1:], "\n")
			break
		}
	}
	err = ioutil.WriteFile(p, []byte(s), os.ModeAppend)
	if err != nil {
		return
	}
	if IsV2RayRunning() {
		err = RestartV2rayService()
	}
	return
}
