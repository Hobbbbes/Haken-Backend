package container

import (
	"fmt"
	"log"

	lxd "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
	"github.com/poodlenoodle42/Hacken-Backend/config"
)

var connection lxd.InstanceServer
var conf config.ContainerConfig

//InstanceChan channel where instances can be recived and returned
var InstanceChan chan string

//InitInstances does initialization of module and starts instances
func InitInstances(co config.ContainerConfig) {
	conf = co
	InstanceChan = make(chan string)
	c, err := lxd.ConnectLXDUnix("", nil)
	if err != nil {
		log.Panic(err)
	}
	connection = c
	for i := 0; i < conf.NumContainer; i++ {
		log.Printf("InitContainer: Start Instance %d\n", i)
		err := startInstance(fmt.Sprintf("Haken%d", i))
		if err != nil {
			log.Panic(err)
		}
		InstanceChan <- fmt.Sprintf("Haken%d", i)
	}
}

func startInstance(name string) error {
	req := api.ContainersPost{
		Name: name,
		Source: api.ContainerSource{
			Type:  "image",
			Alias: "Haken",
		},
	}
	op, err := connection.CreateContainer(req)
	if err != nil {
		return err
	}

	// Wait for the operation to complete
	err = op.Wait()
	if err != nil {
		return err
	}
	i, _, err := connection.GetInstance(name)
	prof := i.InstancePut
	prof.Profiles = append(prof.Profiles, "Haken")
	op, err = connection.UpdateInstance(name, prof, "")
	if err != nil {
		return err
	}
	err = op.Wait()
	if err != nil {
		return err
	}
	// Get LXD to start the container (background operation)
	reqState := api.ContainerStatePut{
		Action:  "start",
		Timeout: 10,
	}

	op, err = connection.UpdateContainerState("h2", reqState, "")
	if err != nil {
		return err
	}
	// Wait for the operation to complete
	err = op.Wait()
	if err != nil {
		return err
	}
	return nil
}

//StopAndDeleteInstances stops and deletes all running Haken container instances
func StopAndDeleteInstances() {
	for i := 0; i < conf.NumContainer; i++ {
		log.Printf("StopAndDeleteInstances: Stop Instance %d\n", i)
		err := startInstance(fmt.Sprintf("Haken%d", i))
		if err != nil {
			log.Panic(err)
		}
	}
}

func stopAndDeleteInstance(name string) error {
	reqState := api.ContainerStatePut{
		Action:  "stop",
		Timeout: 10,
	}

	op, err := connection.UpdateContainerState(name, reqState, "")
	if err != nil {
		return err
	}
	// Wait for the operation to complete
	err = op.Wait()
	if err != nil {
		return err
	}
	op, err = connection.DeleteContainer(name)

	if err != nil {
		return err
	}
	// Wait for the operation to complete
	err = op.Wait()
	if err != nil {
		return err
	}
	return nil
}
