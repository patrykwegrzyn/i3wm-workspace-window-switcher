package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"go.i3wm.org/i3/v4"
)

type Window struct {
	id         int64
	name       string
	focused    bool
	fullscreen int64
}

func (w *Window) isFullscreen() bool {
	return w.fullscreen > 0
}

type Workspace struct {
	name    string
	windows map[string]Window
}

func (ww *Workspace) getWindow(name string) (Window, error) {
	if val, ok := ww.windows[name]; ok {
		return val, nil
	}
	err := fmt.Errorf("%s \n not a valid selection please check your command", name)
	return Window{}, err

}

func (ww *Workspace) getWindowNamesAsBuffer() bytes.Buffer {
	buffer := bytes.Buffer{}
	for _, w := range ww.windows {
		buffer.Write([]byte(fmt.Sprintln(w.name)))
	}
	return buffer
}

func (ww *Workspace) getFocused() Window {
	var window Window
	for _, w := range ww.windows {
		if w.focused {
			window = w
			break
		}
	}

	return window
}

//W global var for workspace
var W = Workspace{"", make(map[string]Window)}

func main() {

	workspace := getActiveWorkspace()
	W.name = workspace.Name
	for _, n := range findChildNodes(workspace.Nodes) {
		W.windows[n.Name] = Window{
			int64(n.ID), n.Name, n.Focused, int64(n.FullscreenMode),
		}
	}

	buffer := W.getWindowNamesAsBuffer()
	selection, err := doComand(buffer, W.name)
	handleError(err)

	focused := W.getFocused()
	fullscreen := focused.isFullscreen()

	window, err := W.getWindow(selection)
	handleError(err)

	if fullscreen {
		_, err := i3Command(window.id, "fullscreen")
		handleError(err)
	}

	_, err2 := i3Command(window.id, "focus")
	handleError(err2)

}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func doComand(b bytes.Buffer, name string) (string, error) {
	mode, flags := parseArgs()
	_, err := exec.LookPath("rofi")

	handleError(err)

	flags = append(flags, []string{"-p", name}...)
	cmd := exec.Command(mode, flags...)
	cmd.Stdin = &b

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("here", err)
		return "", err
	}

	selection := strings.TrimSuffix(string(out), "\n")
	//fmt.Println("selection", selection)
	return selection, nil
}

func i3Command(id int64, cmd string) ([]i3.CommandResult, error) {
	c := fmt.Sprintf("[con_id=\"%d\"] %s", id, cmd)
	res, err := i3.RunCommand(c)
	return res, err
}

func findChildNodes(nodes []*i3.Node) []*i3.Node {
	var n []*i3.Node
	for _, node := range nodes {
		if node.Name == "" {
			n = append(n, findChildNodes(node.Nodes)...)
		} else {
			n = append(n, node)
		}
	}
	return n
}

func getActiveWorkspace() *i3.Node {
	tree, _ := i3.GetTree()
	workspaces, _ := i3.GetWorkspaces()
	var workspace i3.Workspace

	for _, w := range workspaces {
		if w.Focused {
			workspace = w
			break
		}
	}

	return tree.Root.FindChild(func(n *i3.Node) bool {
		return int64(n.ID) == int64(workspace.ID)
	})

}

func parseArgs() (mode string, args []string) {
	var m string
	_args := os.Args[1:]

	for i, a := range _args {
		if a == "-mode" {
			m = _args[i+1]
			_args = append(_args[:i], _args[i+2:]...)
		}
	}

	if m == "" {
		m = "dmenu"
		fmt.Println("-mode flag not provided using default", m)
	}

	return m, _args
}
