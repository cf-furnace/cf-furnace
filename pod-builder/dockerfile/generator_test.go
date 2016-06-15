package dockerfile_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/linsun/pod-builder/dockerfile"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Generator", func() {
	var dropletURI, startCommand string

	BeforeEach(func() {
		dropletURI = "file://foo/bar/baz"
		startCommand = "some start command with spaces"
	})

	Describe("CreateRuntimeDockerfile", func() {
		It("generates a Dockerfile that copies the droplet", func() {
			buf := &bytes.Buffer{}

			err := dockerfile.CreateRuntimeDockerfile(buf, dropletURI, startCommand)
			Expect(err).NotTo(HaveOccurred())

			contents := buf.String()
			Expect(contents).To(Equal("FROM busybox:latest\n" +
				"wget file://foo/bar/baz\n" +
				"COPY droplet.tgz droplet.tgz\n" +
				"CMD some start command with spaces"))
		})
	})

	Describe("WriteRuntimeDockerfile", func() {
		var tempDir, path string

		BeforeEach(func() {
			var err error
			tempDir, err = ioutil.TempDir("", "dockerfile")
			Expect(err).NotTo(HaveOccurred())

			path = filepath.Join(tempDir, "Dockerfile")
		})

		AfterEach(func() {
			os.RemoveAll(tempDir)
		})

		It("write the dockerile to the specified location", func() {
			err := dockerfile.WriteRuntimeDockerfile(path, dropletURI, startCommand)
			Expect(err).NotTo(HaveOccurred())
			Expect(path).To(BeAnExistingFile())

			fileContents, err := ioutil.ReadFile(path)
			Expect(err).NotTo(HaveOccurred())

			expectedContents := &bytes.Buffer{}
			err = dockerfile.CreateRuntimeDockerfile(expectedContents, dropletURI, startCommand)
			Expect(err).NotTo(HaveOccurred())

			Expect(fileContents).To(Equal(expectedContents.Bytes()))
		})

		Context("when write file fails", func() {
			BeforeEach(func() {
				f, err := os.Create(path)
				Expect(err).NotTo(HaveOccurred())
				f.Close()

				err = os.Chmod(path, 0444)
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns an error", func() {
				err := dockerfile.WriteRuntimeDockerfile(path, dropletURI, startCommand)
				Expect(err).To(MatchError(MatchRegexp(".*permission denied")))
			})
		})
	})
})
