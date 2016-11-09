package cmd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/namesgenerator"
	"golang.org/x/net/context"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

func base64EncodeAuth(auth types.AuthConfig) (string, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(auth); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(buf.Bytes()), nil
}

func printContainerLogs(cli *client.Client, resp types.ContainerCreateResponse, ctx context.Context) ([]byte, error) {
	out, err := cli.ContainerLogs(
		ctx,
		resp.ID,
		types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
		})
	if err != nil {
		return nil, err
	}

	defer out.Close()

	content, err := ioutil.ReadAll(out)
	if err != nil {
		return nil, err
	}

	return content, nil
}

// post cluster help types
type helptype int

const (
	Created helptype = iota
	Destroyed
	Updated
)

func clusterHelpError(help helptype, clusterConfigFile string) {
	fmt.Println("Some of the cluster state MAY be available:")
	clusterHelp(help, clusterConfigFile)
}

func clusterHelp(help helptype, clusterConfigFile string) {
	if _, err := os.Stat(path.Join(outputLocation,
		clusterConfig.GetString("deployment.cluster"), "admin.kubeconfig")); err == nil {
		fmt.Println("To use kubectl: ")
		fmt.Println(" kubectl --kubeconfig=" + path.Join(
			outputLocation,
			clusterConfig.GetString("deployment.cluster"), "admin.kubeconfig") + " [kubectl commands]")
		fmt.Println(" or use 'k2cli tool kubectl --config " + clusterConfigFile + " [kubectl commands]'")

		if _, err := os.Stat(path.Join(outputLocation,
			clusterConfig.GetString("deployment.cluster"), "admin.kubeconfig")); err == nil {
			fmt.Println("To use helm: ")
			fmt.Println(" export KUBECONFIG=" + path.Join(
				outputLocation,
				clusterConfig.GetString("deployment.cluster"), "admin.kubeconfig"))
			fmt.Println(" helm [helm command] --home" + path.Join(
				outputLocation,
				clusterConfig.GetString("deployment.cluster"), ".helm"))
			fmt.Println(" or use 'k2cli tool helm --config " + clusterConfigFile + " [helm commands]'")
		}
	}

	if _, err := os.Stat(path.Join(outputLocation,
		clusterConfig.GetString("deployment.cluster"), "ssh_config")); err == nil {
		fmt.Println("To use ssh: ")
		fmt.Println(" ssh <node pool name>-<number> -F " + path.Join(
			outputLocation,
			clusterConfig.GetString("deployment.cluster"), "ssh_config"))
		fmt.Println(" or use 'k2cli tool --config ssh ssh " + clusterConfigFile + " [ssh commands]'")
	}
}

func containerEnvironment() []string {
	envs := []string{
		"ANSIBLE_NOCOLOR=True",
		"DISPLAY_SKIPPED_HOSTS=0",
		"AWS_ACCESS_KEY_ID=" + os.Getenv("AWS_ACCESS_KEY_ID"),
		"AWS_SECRET_ACCESS_KEY=" + os.Getenv("AWS_SECRET_ACCESS_KEY"),
		"AWS_DEFAULT_REGION=" + os.Getenv("AWS_DEFAULT_REGION"),
		"CLOUDSDK_COMPUTE_ZONE=" + os.Getenv("CLOUDSDK_COMPUTE_ZONE"),
		"CLOUDSDK_COMPUTE_REGION=" + os.Getenv("CLOUDSDK_COMPUTE_REGION"),
		"KUBECONFIG=" + path.Join(outputLocation,
			clusterConfig.GetString("deployment.cluster"),
			"admin.kubeconfig"),
		"HELM_HOME=" + path.Join(outputLocation,
			clusterConfig.GetString("deployment.cluster"),
			".helm"),
	}

	return envs
}

func makeMounts(clusterConfigPath string) *container.HostConfig {
	// cluster configuration is always mounted
	var hostConfig *container.HostConfig
	if len(strings.TrimSpace(clusterConfigPath)) > 0 {
		hostConfig = &container.HostConfig{
			Binds: []string{
				clusterConfigPath + ":" + clusterConfigPath,
				outputLocation + ":" + outputLocation},
		}

		deployment := reflect.ValueOf(clusterConfig.Sub("deployment"))
		parseMounts(deployment, hostConfig)
	} else {
		hostConfig = &container.HostConfig{
			Binds: []string{
				outputLocation + ":" + outputLocation},
		}
	}

	return hostConfig
}

func parseMounts(deployment reflect.Value, hostConfig *container.HostConfig) {
	switch deployment.Kind() {
	case reflect.Ptr:
		deploymentValue := deployment.Elem()

		// Check if the pointer is nil
		if !deploymentValue.IsValid() {
			return
		}

		parseMounts(deploymentValue, hostConfig)

	case reflect.Interface:
		deploymentValue := deployment.Elem()
		parseMounts(deploymentValue, hostConfig)

	case reflect.Struct:
		for i := 0; i < deployment.NumField(); i += 1 {
			parseMounts(deployment.Field(i), hostConfig)
		}

	case reflect.Slice:
		for i := 0; i < deployment.Len(); i += 1 {
			parseMounts(deployment.Index(i), hostConfig)
		}

	case reflect.Map:
		for _, key := range deployment.MapKeys() {
			originalValue := deployment.MapIndex(key)
			parseMounts(originalValue, hostConfig)
		}
	case reflect.String:
		reflectedString := fmt.Sprintf("%s", deployment)
		if _, err := os.Stat(os.ExpandEnv(reflectedString)); err == nil {
			if filepath.IsAbs(os.ExpandEnv(reflectedString)) {
				fmt.Println(reflectedString)
				hostConfig.Binds = append(hostConfig.Binds, os.ExpandEnv(reflectedString)+":"+os.ExpandEnv(reflectedString))
			}
		}
	default:
	}
}

func getClient() *client.Client {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient(dockerHost, "", nil, defaultHeaders)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return cli
}

func getAuthConfig64(cli *client.Client, ctx context.Context) string {
	authConfig := types.AuthConfig{}
	if len(userName) > 0 && len(password) > 0 {
		imageParts := strings.Split(containerImage, "/")
		if strings.Count(imageParts[0], ".") > 0 {
			authConfig.ServerAddress = imageParts[0]
		} else {
			authConfig.ServerAddress = "index.docker.io"
		}

		authConfig.Username = userName
		authConfig.Password = password

		_, err := cli.RegistryLogin(ctx, authConfig)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}

	base64Auth, err := base64EncodeAuth(authConfig)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return base64Auth
}

func pullImage(cli *client.Client, ctx context.Context, base64Auth string) {

	pullOpts := types.ImagePullOptions{
		RegistryAuth:  base64Auth,
		All:           false,
		PrivilegeFunc: nil,
	}

	pullResponseBody, err := cli.ImagePull(ctx, containerImage, pullOpts)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	defer pullResponseBody.Close()

	// wait until the image download is finished
	dec := json.NewDecoder(pullResponseBody)
	m := map[string]interface{}{}
	for {
		if err := dec.Decode(&m); err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			panic(err)
		}
	}

	// if the final stream object contained an error, panic
	if errMsg, ok := m["error"]; ok {
		fmt.Println("%v", errMsg)
		panic(errMsg)
	}
}

func containerAction(cli *client.Client, ctx context.Context, command []string, k2config string) (types.ContainerCreateResponse, int) {
	containerConfig := &container.Config{
		Image:        containerImage,
		Env:          containerEnvironment(),
		Cmd:          command,
		AttachStdout: true,
		Tty:          true,
	}

	hostConfig := makeMounts(k2config)
	resp, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, "k2-"+namesgenerator.GetRandomName(1))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		fmt.Println(err)
		panic(err)
	}

	statusCode, err := cli.ContainerWait(ctx, resp.ID)
	if err != nil {
		select {
		case <-ctx.Done():
			fmt.Println("Action timed out!")
			return resp, 1
		default:
			fmt.Println(err)
			panic(err)
		}
	}

	return resp, statusCode
}

func getContext() (ctx context.Context) {
	return context.Background()
}

func getTimedContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(actionTimeout)*time.Second)
}

func writeLog(logFilePath string, out []byte) {
	var fileHandle *os.File

	_, err := os.Stat(logFilePath)
	if err != nil {
		if os.IsNotExist(err) {

			// make sure path exists
			err = os.MkdirAll(filepath.Dir(logFilePath), 0777)
			if err != nil {
				fmt.Println(err)
				panic(err)
			}

			// check if a valid file path
			var d []byte
			if err := ioutil.WriteFile(logFilePath, d, 0644); err == nil {
				os.Remove(logFilePath)
			} else {
				fmt.Println(err)
				panic(err)
			}

			fileHandle, err = os.Create(logFilePath)
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
		} else {
			fileHandle, err = os.OpenFile("test.txt", os.O_RDWR, 0666)
		}
	}

	defer fileHandle.Close()

	_, err = fileHandle.Write(out)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
