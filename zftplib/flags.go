package zftplib

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

type FtpFlags struct {
	ZConfig     string
	Debug       bool
	ServiceName string
}



func (f *FtpFlags) GetUserAndIdentity(input string) (string, string) {
	username := ParseUserName(input)
	f.DebugLog("      username set to: %s", username)
	targetIdentity := ParseTargetIdentity(input)
	f.DebugLog("targetIdentity set to: %s", targetIdentity)
	return username, targetIdentity
}

func ParseUserName(input string) string {
	var username string
	if strings.ContainsAny(input, "@") {
		userServiceName := strings.Split(input, "@")
		username = userServiceName[0]
	} else {
		curUser, err := user.Current()
		if err != nil {
			logrus.Fatal(err)
		}
		username = curUser.Username
		if strings.Contains(username, "\\") && runtime.GOOS == "windows" {
			username = strings.Split(username, "\\")[1]
		}
	}
	return username
}

func ParseTargetIdentity(input string) string {
	var targetIdentity string
	if strings.ContainsAny(input, "@") {
		targetIdentity = strings.Split(input, "@")[1]
	} else {
		targetIdentity = input
	}

	if strings.Contains(targetIdentity, ":") {
		return strings.Split(targetIdentity, ":")[0]
	}
	return targetIdentity
}

func ParseFilePath(input string) string {
	if strings.Contains(input, ":") {
		colPos := strings.Index(input, ":") + 1
		return input[colPos:]
	}
	return input
}

func (f *FtpFlags) InitFlags(cmd *cobra.Command, exeName string) {
	cmd.PersistentFlags().StringVarP(&f.ServiceName, "service", "s", exeName, fmt.Sprintf("service name. default: %s", exeName))
	cmd.PersistentFlags().StringVarP(&f.ZConfig, "ZConfig", "c", "", fmt.Sprintf("Path to ziti config file. default: $HOME/.ziti/%s.json", f.ServiceName))
	cmd.PersistentFlags().BoolVarP(&f.Debug, "debug", "d", false, "pass to enable additional debug information")

	if f.ZConfig == "" {
		userHome, err := os.UserHomeDir()
		if err != nil {
			logrus.Fatalf("could not find UserHomeDir? %v", err)
		}
		f.ZConfig = filepath.Join(userHome, ".ziti", fmt.Sprintf("%s.json", exeName))
	}
	f.DebugLog("       ZConfig set to: %s", f.ZConfig)
}
