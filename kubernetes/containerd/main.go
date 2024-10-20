package main

import (
	"context"
	"log"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
)

func containerExample() error {
	client, err := containerd.New("/run/k3s/containerd/containerd.sock")
	if err != nil {
		return err
	}

	ctx := namespaces.WithNamespace(context.Background(), "default")
	image, err := client.Pull(ctx, "docker.io/library/busybox:1.36", containerd.WithPullUnpack)
	if err != nil {
		return err
	}

	log.Printf("Successfully pulled %s image\n", image.Name())

	container, err := client.NewContainer(
		ctx,
		"busybox",
		containerd.WithNewSnapshot("busybox", image),
		containerd.WithNewSpec(oci.WithImageConfig(image), oci.WithProcessArgs("sleep", "infinity")),
	)
	if err != nil {
		return err
	}

	defer container.Delete(ctx, containerd.WithSnapshotCleanup)
	log.Printf("Successfully created container with ID %s and snapshot with ID busybox", container.ID())

	return nil
}

func main() {
	if err := containerExample(); err != nil {
		log.Fatal(err)
	}
}
