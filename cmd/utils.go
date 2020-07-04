package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yottab/cli/config"
	ybApi "github.com/yottab/proto-api/proto"
	"golang.org/x/crypto/ssh/terminal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	flagVariableArray = make([]string, 0, 8)
	flagIndex         int32
	flagAppName       string
	flagFile          string
	clientTimeout     time.Duration = 10 * time.Second
)

func parseDotEnvFile(filename string) map[string]string {
	var myEnv map[string]string
	myEnv, err := godotenv.Read(filename)
	if err != nil {
		panic(err)
	}
	return myEnv
}
func arrayFlagToMap(flags []string) map[string]string {
	varMap := make(map[string]string, len(flags))
	for _, v := range flags {
		index := strings.Index(v, "=")
		if index > 0 {
			key := v[:index]
			varMap[key] = v[index+1:]
		} else {
			log.Fatalf("Bad data entry format [%s], Enter the information in 'KEY=VALUE' format.", v)
		}
	}
	return varMap
}

func forceFlagGetStrValue(cmd *cobra.Command, flagName, inputAnswr string) (val string) {
	val = cmd.Flag(flagName).Value.String()
	if val == "" {
		val = readFromConsole(inputAnswr)
	}
	return
}

func forceArgGetStrValue(args []string, index int, inputAnswr string) (val string) {
	val = getCliRequiredArg(args, index)
	if val == "" {
		val = readFromConsole(inputAnswr)
	}
	return
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

func grpcConnect() ybApi.Client {
	return ybApi.Connect(
		viper.GetString(config.KEY_HOST),
		ybApi.NewPerRPC(func() string {
			return viper.GetString(config.KEY_TOKEN)
		}, func() string {
			return version
		}, nil), clientTimeout)
}

func toTime(t *ybApi.Timestamp) (out string) {
	if t != nil {
		out = time.Unix(t.Seconds, 0).Format(time.RFC3339)
	}
	return
}

func endpointTypeValid(etype string) error {
	switch etype {
	case "http", "grpc", "private":
		return nil
	default:
		return errors.New("Endpoint type is invalid, valid values are http, grpc and private")
	}
}

func streamAppLog(args []string) {
	var client ybApi.Client
	firstTry := true
	for {
		if !firstTry {
			//Wait and retry
			time.Sleep(time.Millisecond * 500)
		}
		client = grpcConnect()
		req := getCliRequestIdentity(args, 0)
		logClient, err := client.V2().AppLog(context.Background(), req)
		uiCheckErr("Could not Get Application log", err)
		err = uiStreamLog(logClient)
		client.Close()
		if err != nil {
			if status.Code(err) == codes.ResourceExhausted {
				break
			}
			if strings.Contains(err.Error(), "RST_STREAM") {
				//Resume log streaming on proto related error
				log.Println("Timeout streaming log...")
				// continue
				return
			}
		} else {
			break
		}

	}

}

func streamBuildLog(appName, appTag string, waitToReady bool) {
	var client ybApi.Client
	firstTry := true
	id := getRequestIdentity(
		fmt.Sprintf(pushLogIDFormat, appName, appTag))
	for {
		if !firstTry {
			//Wait and retry
			time.Sleep(time.Millisecond * 500)
		}
		client = grpcConnect()
		logClient, err := client.V2().ImgBuildLog(context.Background(), id)
		if err != nil && waitToReady {
			if status.Code(err) == codes.NotFound {
				log.Println("Log not ready yet, we try again in 20 seconds...")
				time.Sleep(time.Second * 20)
				continue
			}
		}
		uiCheckErr(fmt.Sprintf("Could not get build log right now!\nTry again in a few soconds using:\n$yb push log"), err)
		err = uiStreamLog(logClient)

		client.Close()
		if err != nil {
			if status.Code(err) == codes.ResourceExhausted {
				break
			}
			if strings.Contains(err.Error(), "RST_STREAM") {
				//Resume log streaming on proto related error
				log.Println("Timeout streaming log...")
				// continue
				return
			}
		} else {
			break
		}
	}
}

//wrapRemove is a middleware for confirm object removes
func wrapRemove(objectType string, f func(cmd *cobra.Command, args []string)) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		objectName := getCliRequiredArg(args, 0)
		confirmF(cmd, "to remove \033[1m%s\033[0m of type \033[1m%s\033[0m", objectName, objectType)
		//Run the original func if confirmed
		f(cmd, args)
	}
}

//confirmF check for confirmations, exit on no confirmation, otherwise a noop function
func confirmF(cmd *cobra.Command, format string, args ...interface{}) {
	if flagConfirm {
		return
	}
	confirmQ := fmt.Sprintf("Do you confrim %s, yes/no(default: no)?", fmt.Sprintf(format, args...))
	for {
		val := readFromConsole(confirmQ)
		lowerAns := strings.ToLower(val)
		switch lowerAns {
		case "y", "yes":
			return
		case "no", "n", "":
			os.Exit(1)
		default:
			continue
		}
	}
}
