package main
import (
    "os"
    "fmt"
    "time"
    "errors"
    "strings"

	docker "github.com/fsouza/go-dockerclient"
    pushbullet "github.com/xconstruct/go-pushbullet"
)

type DockerMonitor struct {
    Docker *docker.Client
    Conts []DockerContainer
    Pb *pushbullet.Client
    PbDevs []*pushbullet.Device
}

type DockerContainer struct {
    Id string
    Name string
    Healthy int
}

func main() {
    dm, err := New(os.Getenv("DOCKERMONITOR_DOCKERSOCK"), os.Getenv("DOCKERMONITOR_PBTOKEN"))

    if err != nil {
        panic(err)
    }

    for true {
        fmt.Println("Heartbeat!")
        dm.Heartbeat()
        time.Sleep(600 * time.Second)
    }
}

func New(dc string, pbToken string) (DockerMonitor, error) {
    var dockerMonitor DockerMonitor

    // pushbullet
    pb := pushbullet.New(pbToken)
    devs, err := pb.Devices()
    if err != nil {
        return dockerMonitor, err
    }

    dockerClient, err := docker.NewClient(dc)
    if err != nil {
        return dockerMonitor, err
    }

    dockerMonitor = DockerMonitor{
        Docker: dockerClient,
        Pb: pb,
        PbDevs: devs,
    }

    return dockerMonitor, nil
}

func (dm *DockerMonitor) Heartbeat() error {
    conts, err := dm.getContainers()
    if err != nil {
        return err
    }

    err = dm.checkUnhealthy(conts)
    if err != nil {
        return err
    }

    return nil
}

func (dm *DockerMonitor) getContainers() ([]DockerContainer, error) {
    containers, err := dm.Docker.ListContainers(docker.ListContainersOptions{})

    if err != nil {
        return nil, errors.New("could not get containers.")
    }

    var conts []DockerContainer
    var healthy int

    for _, cont := range containers {
        if strings.Contains(cont.Status, "unhealthy") {
            healthy = 0
        } else if strings.Contains(cont.Status, "healthy") {
            healthy = 1
        } else {
            healthy = 2
        }


        conts = append(conts, DockerContainer{
            Id: cont.ID,
            Name: strings.Join(cont.Names[:], ","),
            Healthy: healthy,
        })
    }

    return conts, nil
}

func (dm *DockerMonitor) checkUnhealthy(conts []DockerContainer) error {
    var err error
    for _, cont := range conts {
        if cont.Healthy != 0 {
            continue
        }
        err = dm.sendNotification("Docker Alert!", "Container - "+cont.Name+" is unhealthy")

        if err != nil {
            return err
        }
    }

    return nil
}

func (dm *DockerMonitor) sendNotification(title string, msg string) error {
    return dm.Pb.PushNote(dm.PbDevs[0].Iden, title, msg)
}
