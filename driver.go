package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/docker/go-plugins-helpers/volume"
)

type ARVol struct {
	volumes    map[string]string
	m          *sync.Mutex
	mountPoint string
}

func NewARVol() ARVol {
	return ARVol{
		volumes:    make(map[string]string),
		m:          &sync.Mutex{},
		mountPoint: "/var/lib/ARVol/mount",
	}
}

func (d ARVol) Create(r volume.Request) volume.Response {
	logrus.Infof("Create volume: %s", r.Name)
	d.m.Lock()
	defer d.m.Unlock()

	if _, ok := d.volumes[r.Name]; ok {
		return volume.Response{}
	}

	volumePath := filepath.Join(d.mountPoint, r.Name)

	_, err := os.Lstat(volumePath)
	if err != nil {
		logrus.Errorf("Error %s %v", volumePath, err.Error())
		return volume.Response{Err: fmt.Sprintf("Error: %s: %s", volumePath, err.Error())}
	}

	d.volumes[r.Name] = volumePath

	return volume.Response{}
}

func (d ARVol) List(r volume.Request) volume.Response {
	logrus.Info("Volumes list ", r)

	volumes := []*volume.Volume{}

	for name, path := range d.volumes {
		volumes = append(volumes, &volume.Volume{
			Name:       name,
			Mountpoint: path,
		})
	}

	return volume.Response{Volumes: volumes}

}

func (d ARVol) Get(r volume.Request) volume.Response {
	logrus.Info("Get volume ", r)
	if path, ok := d.volumes[r.Name]; ok {
		return volume.Response{
			Volume: &volume.Volume{
				Name:       r.Name,
				Mountpoint: path,
			},
		}
	}
	return volume.Response{
		Err: fmt.Sprintf("volume named %s not found", r.Name),
	}
}

func (d ARVol) Remove(r volume.Request) volume.Response {
	logrus.Info("Remove volume ", r)

	d.m.Lock()
	defer d.m.Unlock()

	if _, ok := d.volumes[r.Name]; ok {
		delete(d.volumes, r.Name)
	}

	return volume.Response{}
}

func (d ARVol) Path(r volume.Request) volume.Response {
	logrus.Info("Get volume path", r)

	if path, ok := d.volumes[r.Name]; ok {
		return volume.Response{
			Mountpoint: path,
		}
	}
	return volume.Response{}
}

func (d ARVol) Capabilities(r volume.Request) volume.Response {

	return volume.Response{Capabilities: volume.Capability{Scope: "global"}}
}

//func (d ARVol) Mount(r volume.Request) volume.Response {
func (d ARVol) Mount(r volume.MountRequest) volume.Response {
	logrus.Info("Mount volume ", r)

	if path, ok := d.volumes[r.Name]; ok {
		return volume.Response{
			Mountpoint: path,
		}
	}

	return volume.Response{}

}

//func (d ARVol) Unmount(r volume.Request) volume.Response {
func (d ARVol) Unmount(r volume.UnmountRequest) volume.Response {
	logrus.Info("Unmount ", r)
	return volume.Response{}
}
