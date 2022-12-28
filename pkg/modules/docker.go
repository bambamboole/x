package modules

import (
	"bytes"
	"text/template"
)

type DockerConfig struct {
	Commands map[string]struct {
		Name    string
		Service string
		Bin     string
	}
}

type DockerModule struct {
	config *DockerConfig
}

func (d *DockerModule) GetConfig() interface{} {
	if d.config == nil {
		d.config = &DockerConfig{}
	}
	return d.config
}

func (d *DockerModule) GetScript() (string, error) {
	var b bytes.Buffer

	for name, cmd := range d.config.Commands {
		t, err := template.New("command").Parse(`
# docker-compose options

docker:container_run() {
	# Parse container name
	local container="{{"${1:-}"}}"
	shift
	if [ -z "$container" ]
	then
		x:error "You need to pass the container name"
		return 1
	fi
	(
		cd $X_BASE_PATH || return 1
		docker compose run "$container" $@
	)
}

docker:{{.Name}}() {
    docker:container_run {{.Service}} "{{.Bin}} $@"
}

`)
		if err != nil {
			return b.String(), err
		}
		err = t.Execute(&b, struct {
			Name    string
			Service string
			Bin     string
		}{Name: name, Service: cmd.Service, Bin: cmd.Bin})
		if err != nil {
			return b.String(), err
		}
	}
	return b.String(), nil
}
