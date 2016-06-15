package dockerfile

import (
	"fmt"
	"io"
	"io/ioutil"
)

func CreateRuntimeDockerfile(writer io.Writer, dropletURI string, startCommand string) error {
	fmt.Fprintf(writer, appDataContainerInitialDockerfileString(dropletURI, startCommand))
	return nil
}

func WriteRuntimeDockerfile(path, dropletUri, startCmd string) error {
	dockerfile := appDataContainerInitialDockerfileString(dropletUri, startCmd)
	return ioutil.WriteFile(path, []byte(dockerfile), 0644)
}

func appDataContainerInitialDockerfileString(dropletUri string, startCmd string) string {
	return `FROM busybox:latest
wget ` + dropletUri + `
COPY droplet.tgz droplet.tgz
CMD ` + startCmd
}
