package transformer_test

import (
	"github.com/cloudfoundry-incubator/bbs/models"
	"github.com/cloudfoundry-incubator/runtime-schema/cc_messages"
	"github.com/linsun/pod-builder/transformer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/kubernetes/pkg/api"
)

var _ = Describe("Transformer", func() {
	var desiredApp cc_messages.DesireAppRequestFromCC
	var expectedPod *api.Pod

	BeforeEach(func() {
		desiredApp = cc_messages.DesireAppRequestFromCC{
			ProcessGuid:  "process-guid-1",
			DropletUri:   "source-url-1",
			Stack:        "stack-1",
			StartCommand: "start-command-1",
			Environment: []*models.EnvironmentVariable{
				{Name: "env-key-1", Value: "env-value-1"},
				{Name: "env-key-2", Value: "env-value-2"},
			},
			MemoryMB:        256,
			DiskMB:          1024,
			FileDescriptors: 16,
			NumInstances:    2,
			LogGuid:         "log-guid-1",
			ETag:            "1234567.1890",
			Ports:           []uint32{8080},
		}

		expectedPod = &api.Pod{
			ObjectMeta: api.ObjectMeta{
				Name: "process-guid-1",
			},
			Spec: api.PodSpec{
				Containers: []api.Container{{
					Name:  "process-guid-1-data",
					Image: "localhost:5000/linsun/process-guid-1-data:latest",
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
					Name:  "process-guid-1-runner",
					Image: "localhost:5000/default/k8s-runner:latest",
					Env: []api.EnvVar{
						{Name: "STARTCMD", Value: "start-command-1"},
						{Name: "ENVVARS", Value: "env-key-1=env-value-1,env-key-2=env-value-2"},
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
	})

	It("generates the expected kubernetes pod struct", func() {
		pod, err := transformer.DesiredAppToPod(desiredApp)
		Expect(err).NotTo(HaveOccurred())
		Expect(pod).To(Equal(expectedPod))
	})
})
