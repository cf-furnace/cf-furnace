package transformer

import (
	"fmt"
	"strings"

	"github.com/cloudfoundry-incubator/runtime-schema/cc_messages"
	"k8s.io/kubernetes/pkg/api"
)

func DesiredAppToPod(desiredApp cc_messages.DesireAppRequestFromCC) (*api.Pod, error) {
	var env string
	for _, envVar := range desiredApp.Environment {
		env = strings.TrimPrefix(fmt.Sprintf("%s,%s=%s", env, envVar.Name, envVar.Value), ",")
	}

	pod := &api.Pod{
		ObjectMeta: api.ObjectMeta{
			Name: desiredApp.ProcessGuid,
		},
		Spec: api.PodSpec{
			Containers: []api.Container{{
				Name:  fmt.Sprintf("%s-data", desiredApp.ProcessGuid),
				Image: fmt.Sprintf("localhost:5000/linsun/%s-data:latest", desiredApp.ProcessGuid),
				Lifecycle: &api.Lifecycle{
					PostStart: &api.Handler{
						Exec: &api.ExecAction{
							Command: []string{"cp", "/droplet.tgz", "/app"},
						},
					},
				},
				VolumeMounts: []api.VolumeMount{{
					Name:      "app-volume",
					MountPath: "/app",
				}},
			}, {
				Name:  fmt.Sprintf("%s-runner", desiredApp.ProcessGuid),
				Image: "localhost:5000/default/k8s-runner:latest",
				Env: []api.EnvVar{
					{Name: "STARTCMD", Value: desiredApp.StartCommand},
					{Name: "ENVVARS", Value: env},
					{Name: "PORT", Value: "8080"},
				},
				VolumeMounts: []api.VolumeMount{{
					Name:      "app-volume",
					MountPath: "/app/droplet",
				}},
			}},
			Volumes: []api.Volume{{
				Name: "app-volume",
				VolumeSource: api.VolumeSource{
					EmptyDir: &api.EmptyDirVolumeSource{},
				},
			}},
		},
	}

	return pod, nil
}
