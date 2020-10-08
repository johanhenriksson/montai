package cli

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const imageName = "johanhenriksson/montai"

func IsContainer() bool {
	path := "/proc/1/cgroup"

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return false
	}

	return strings.Contains(string(content), "docker")
}

func Run() {
	fmt.Println("montai cli")

	root := "/mnt/montai"

	// setup
	b, _ := exec.Command(
		"docker", "run", "--rm", "--privileged", "--pid=host",
		imageName, "sh", "install.sh", root).Output()
	fmt.Println(string(b))

	docker, err := client.NewEnvClient()
	if err != nil {
		fmt.Println("Unable to create docker client")
		panic(err)
	}

	containerPort, err := nat.NewPort("tcp", "8080")
	if err != nil {
		log.Fatal("Unable to get the port")
	}

	containerConfig := container.Config{
		Image: imageName,
		Env: []string{
			fmt.Sprintf("ROOT_DIR=%s", root),
		},
		ExposedPorts: nat.PortSet{
			"8080/tcp": struct{}{},
		},
	}
	hostConfig := container.HostConfig{
		AutoRemove: true,
		PortBindings: nat.PortMap{
			containerPort: []nat.PortBinding{
				{HostIP: "0.0.0.0", HostPort: "8080"},
			},
		},
		Mounts: []mount.Mount{
			{
				Type:     "bind",
				Source:   root,
				Target:   root,
				ReadOnly: false,
				BindOptions: &mount.BindOptions{
					Propagation: mount.PropagationRShared,
				},
			},
		},
		Resources: container.Resources{
			Devices: []container.DeviceMapping{
				{
					PathInContainer:   "/dev/fuse",
					PathOnHost:        "/dev/fuse",
					CgroupPermissions: "rwm",
				},
			},
		},
		Privileged: true,
	}

	cont, err := docker.ContainerCreate(context.Background(), &containerConfig, &hostConfig, nil, "montai")
	if err != nil {
		log.Fatal(err)
	}

	err = docker.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})
	if err != nil {
		log.Fatal(err)
	}

	i, err := docker.ContainerLogs(context.Background(), cont.ID, types.ContainerLogsOptions{
		ShowStderr: true,
		ShowStdout: true,
		Timestamps: false,
		Follow:     true,
		Tail:       "40",
	})
	if err != nil {
		log.Fatal(err)
	}
	hdr := make([]byte, 8)
	for {
		_, err := i.Read(hdr)
		if err != nil {
			log.Fatal(err)
		}
		var w io.Writer
		switch hdr[0] {
		case 1:
			w = os.Stdout
		default:
			w = os.Stderr
		}
		count := binary.BigEndian.Uint32(hdr[4:])
		dat := make([]byte, count)
		_, err = i.Read(dat)
		fmt.Fprint(w, string(dat))
	}
}
