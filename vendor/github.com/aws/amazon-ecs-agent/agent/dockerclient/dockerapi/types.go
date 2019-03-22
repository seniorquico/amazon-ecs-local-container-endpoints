// Copyright 2014-2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package dockerapi

import (
	"fmt"
	"time"

	apicontainer "github.com/aws/amazon-ecs-agent/agent/api/container"
	apicontainerstatus "github.com/aws/amazon-ecs-agent/agent/api/container/status"
	apierrors "github.com/aws/amazon-ecs-agent/agent/api/errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/docker/docker/api/types"
)

// ContainerNotFound is a type for a missing container
type ContainerNotFound struct {
	// TaskArn is the ARN of the task the container belongs to
	TaskArn string
	// ContainerName is the name of the container that's missing
	ContainerName string
}

// Error returns an error string for the ContainerNotFound error
func (cnferror ContainerNotFound) Error() string {
	return fmt.Sprintf("Could not find container '%s' in task '%s'",
		cnferror.ContainerName, cnferror.TaskArn)
}

// DockerContainerChangeEvent is a type for container change events
type DockerContainerChangeEvent struct {
	// Status represents the container's status in the event
	Status apicontainerstatus.ContainerStatus
	// DockerContainerMetadata is the metadata of the container in the event
	DockerContainerMetadata
	// Type is the event type received from docker events
	Type apicontainer.DockerEventType
}

// DockerContainerMetadata is a type for metadata about Docker containers
type DockerContainerMetadata struct {
	// DockerID is the contianer's id generated by Docker
	DockerID string
	// ExitCode contains container's exit code if it has stopped
	ExitCode *int
	// PortBindings is the list of port binding information of the container
	PortBindings []apicontainer.PortBinding
	// Error wraps various container transition errors and is set if engine
	// is unable to perform any of the required container transitions
	Error apierrors.NamedError
	// Volumes contains volume informaton for the container
	Volumes []types.MountPoint
	// Labels contains labels set for the container
	Labels map[string]string
	// CreatedAt is the timestamp of container creation
	CreatedAt time.Time
	// StartedAt is the timestamp of container start
	StartedAt time.Time
	// FinishedAt is the timestamp of container stop
	FinishedAt time.Time
	// Health contains the result of a container health check
	Health apicontainer.HealthStatus
	// NetworkMode denotes the network mode in which the container is started
	NetworkMode string
	// NetworksUnsafe denotes the Docker Network Settings in the container
	NetworkSettings *types.NetworkSettings
}

// ListContainersResponse encapsulates the response from the docker client for the
// ListContainers call.
type ListContainersResponse struct {
	// DockerIDs is the list of container IDs from the ListContainers call
	DockerIDs []string
	// Error contains any error returned when listing containers
	Error error
}

// ListImagesResponse encapsulates the response from the docker client for the
// ListImages call.
type ListImagesResponse struct {
	// ImagesIDs is the list of Images IDs from the ListImages call
	ImageIDs []string
	// RepoTags is the list of Images names from the ListImages call
	RepoTags []string
	// Error contains any error returned when listing images
	Error error
}

// VolumeResponse wrapper for CreateVolume and InspectVolume
// TODO Remove type when migration is complete
type VolumeResponse struct {
	DockerVolume *types.Volume
	Error        error
}

// VolumeResponse wrapper for CreateVolume for SDK Clients
type SDKVolumeResponse struct {
	DockerVolume *types.Volume
	Error        error
}

// ListPluginsResponse is a wrapper for ListPlugins api
type ListPluginsResponse struct {
	Plugins []*types.Plugin
	Error   error
}

// String returns a human readable string of the container change event
func (event *DockerContainerChangeEvent) String() string {
	res := fmt.Sprintf("Status: %s, DockerID: %s", event.Status.String(), event.DockerID)
	res += ", health: " + event.Health.Status.String()

	if event.ExitCode != nil {
		res += fmt.Sprintf(", ExitCode: %d", aws.IntValue(event.ExitCode))
	}

	if len(event.PortBindings) != 0 {
		res += fmt.Sprintf(", PortBindings: %v", event.PortBindings)
	}

	if event.Error != nil {
		res += ", Error: " + event.Error.Error()
	}

	if len(event.Volumes) != 0 {
		res += fmt.Sprintf(", Volumes: %v", event.Volumes)
	}

	if len(event.Labels) != 0 {
		res += fmt.Sprintf(", Labels: %v", event.Labels)
	}

	if !event.CreatedAt.IsZero() {
		res += ", CreatedAt: " + event.CreatedAt.String()
	}
	if !event.StartedAt.IsZero() {
		res += ", StartedAt: " + event.StartedAt.String()
	}
	if !event.FinishedAt.IsZero() {
		res += ", FinishedAt: " + event.FinishedAt.String()
	}

	return res
}