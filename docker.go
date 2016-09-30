package testpack

import (
	"context"
	"github.com/docker/engine-api/client"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types"
)

type Container struct {
	ctx    context.Context
	cancel context.CancelFunc
	cli    *client.Client
	id     string

	Image string

	// container -> host
	Ports map[string]string
}

func NewContainer(ctx context.Context, image string, cmd []string, port ...string) (s *Container, err error) {
	s = &Container{
		Image: image,
		Ports: map[string]string{},
	}
	s.ctx, s.cancel = context.WithCancel(ctx)
	s.cli, err = client.NewClient("unix:///var/run/docker.sock", "", nil, nil)
	if err != nil {
		return
	}

	hostPorts, err := GetOpenPorts(len(port))
	if err != nil {
		return
	}
	var portSpecs []string
	for i, iPort := range port {
		hostPort := hostPorts[i]
		s.Ports[iPort] = fmt.Sprintf("%d", hostPort)
		portSpecs = append(portSpecs, fmt.Sprintf("%d:%s", hostPort, iPort))
	}
	_, portBindings, err := nat.ParsePortSpecs(portSpecs)
	if err != nil {
		return
	}
	containerConfig := &container.Config{
		Image: image,
		Cmd: cmd,
		ExposedPorts: map[nat.Port]struct{}{},
	}
	for natPort := range portBindings {
		containerConfig.ExposedPorts[natPort] = struct {}{}
	}

	resp, err := s.cli.ContainerCreate(
		s.ctx,
		containerConfig,
		&container.HostConfig{PortBindings: portBindings},
		nil,
		"",
	)
	if err != nil {
		s.cancel()
		return
	}
	s.id = resp.ID
	err = s.cli.ContainerStart(s.ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		s.cancel()
		return
	}

	return
}

func (d *Container) Pause() (err error) {
	err = d.cli.ContainerPause(d.ctx, d.id)
	return
}

func (d *Container) Unpause() (err error) {
	err = d.cli.ContainerUnpause(d.ctx, d.id)
	return
}

func (d *Container) Close() (err error) {
	d.cli.ContainerStop(d.ctx, d.id, nil)
	d.cli.ContainerRemove(d.ctx, d.id, types.ContainerRemoveOptions{
		Force:true,
	})
	d.cancel()
	return
}

