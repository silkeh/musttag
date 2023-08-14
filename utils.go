package musttag

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

func getMainModules() (modules []string, err error) {
	args := []string{"go", "list", "-m", "-json"}

	data, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
		return nil, fmt.Errorf("running `%s`: %w", strings.Join(args, " "), err)
	}

	var module struct {
		Path      string `json:"Path"`
		Main      bool   `json:"Main"`
		Dir       string `json:"Dir"`
		GoMod     string `json:"GoMod"`
		GoVersion string `json:"GoVersion"`
	}

	decoder := json.NewDecoder(bytes.NewBuffer(data))

	for {
		if err := decoder.Decode(&module); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, fmt.Errorf("decoding json: %w: %s", err, string(data))
		}

		if module.Main {
			modules = append(modules, module.Path)
		}
	}

	return
}
