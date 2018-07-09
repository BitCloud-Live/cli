package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/viper"
	uvApi "github.com/uvcloud/uv-api-go/proto"
	"github.com/uvcloud/uv-cli/config"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	flagVariableArray = make([]string, 0, 8)
	flagIndex         int32
)

func arrayFlagToMap(flags []string) map[string]string {
	varMap := make(map[string]string, len(flags))
	for _, v := range flags {
		splitFlag := strings.Split(v, "=")
		if len(splitFlag) >= 2 {
			varMap[splitFlag[0]] = splitFlag[1]
		} else {
			varMap[v] = ""
		}
	}
	return varMap
}

func readFromConsole(inputAnswr string) (val string) {
	fmt.Print(inputAnswr)
	reader := bufio.NewReader(os.Stdin)
	val, _ = reader.ReadString('\n')
	return strings.TrimSpace(val)
}

func readFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func readPasswordFromConsole(inputAnswr string) (val string) {
	fmt.Print(inputAnswr)
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return ""
	}
	password := string(bytePassword)
	return strings.TrimSpace(password)
}

func grpcConnect() uvApi.Client {
	return uvApi.Connect(viper.GetString(config.KEY_LINK), uvApi.NewJwtAccess(func() string {
		return viper.GetString(config.KEY_TOKEN)
	}))
}
