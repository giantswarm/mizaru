package main

import (
	"github.com/fsouza/go-dockerclient"
	"github.com/juju/errgo"
)

type Container struct {
	IP string
}

type DockerInspector struct {
	client *docker.Client
}

func (di DockerInspector) getContainers() ([]Container, error) {
	var cs []Container

	dockerContainers, err := di.client.ListContainers(docker.ListContainersOptions{})
	if err != nil {
		return nil, errgo.Mask(err)
	}

	for _, dockerContainer := range dockerContainers {
		if containerInfo, err := di.client.InspectContainer(dockerContainer.ID); err != nil {
			return nil, errgo.Mask(err)
		} else {
			if containerInfo.NetworkSettings != nil && containerInfo.NetworkSettings.IPAddress != "" {
				c := Container{
					IP: containerInfo.NetworkSettings.IPAddress,
				}

				cs = append(cs, c)
			}
		}
	}

	return cs, nil
}
