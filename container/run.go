package container

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	lxd "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
	"github.com/poodlenoodle42/Hacken-Backend/datastructures"
)

//Exitcodes
const (
	//OK Program ran with no errors
	OK = 0
	//MinorProblem e.g. when compilation fails
	MinorProblem = 1
	TimeOut      = 124
	//CommandNotFound
	CommandNotFound = 127
	//SIGILL illegal instruction or corrupt binary
	SIGILL = 132
	//SIGTRAP Program was aborted perhabs as result of dividing an interger by zero
	SIGTRAP = 133
	//SIGABRT Program was aborted perhabs as result of a failed assertion
	SIGABRT = 134
	//SIGFPE Program was aborted perhabs as result of floating point exception or integer overflow
	SIGFPE      = 136
	OutOfMemory = 137
	//SIGBUS Program was aborted perhabs as result of unaligned memory access
	SIGBUS   = 138
	SegFault = 139
)

//PrepareExecution prepares the execution enviroment by copying source files and executing the PreLaunchTask
func PrepareExecution(SourcePath string, lang datastructures.Language, instance string) (datastructures.Status, error) {
	var s datastructures.Status
	f, err := os.Open(SourcePath)
	if err != nil {
		return s, err
	}
	args := lxd.InstanceFileArgs{
		Content: f,
		UID:     1500,
		GID:     1500,
		Mode:    0o777,
	}
	err = connection.CreateInstanceFile(instance, "/home/runner/main."+lang.Abbreviation, args)
	if err != nil {
		return s, err
	}
	if lang.PreLaunchTask == "" {
		s.ExitCode = -1
		s.Output = ""
		return s, nil
	}
	req := api.ContainerExecPost{
		Command:     strings.Split(lang.PreLaunchTask, " "),
		WaitForWS:   true,
		Interactive: false,
		User:        1500,
		Group:       1500,
		Cwd:         "/home/runner/",
		Environment: map[string]string{"GOPATH": "/home/runner/go", "HOME": "/home/runner"},
	}
	read, write := io.Pipe()
	args2 := lxd.ContainerExecArgs{
		Stdin:  nil,
		Stdout: write,
		Stderr: write,
	}
	reader := bufio.NewReader(read)
	buf := make([]byte, 0, 1024*10)
	op, err := connection.ExecContainer(instance, req, &args2)
	if err != nil {
		return s, err
	}
	// Wait for it to complete
	err = op.Wait()
	if err != nil {
		return s, err
	}
	s.ExitCode = int(op.Get().Metadata["return"].(float64))
	if s.ExitCode == 0 { //Compiled, no error message
		s.ExitCode = -1
		return s, nil
	}
	go func() {
		time.Sleep(time.Millisecond * 1000)
		read.Close()
		write.Close()
	}()
	n, err := reader.Read(buf[:cap(buf)])
	if err != nil {
		if err == io.EOF {
			return s, nil
		}
		return s, err
	}
	s.Output = string(buf[:n])
	return s, nil
}

var timeout = 10
var killTimeout = 2

//Exec executes given task, kills it after time and redirects input, returns output
func Exec(instance string, command string, in io.ReadCloser) (datastructures.Status, error) {
	s := datastructures.Status{
		ExitCode: -1,
		Output:   "",
	}
	c := []string{"/usr/bin/timeout", "-k", strconv.Itoa(killTimeout), strconv.Itoa(timeout)}
	c = append(c, strings.Split(command, " ")...)
	req := api.ContainerExecPost{
		Command:     c,
		WaitForWS:   true,
		Interactive: false,
		User:        1500,
		Group:       1500,
		Cwd:         "/home/runner/",
		Environment: map[string]string{"GOPATH": "/home/runner/go"},
	}
	read, write := io.Pipe()
	args := lxd.ContainerExecArgs{
		Stdin:  in,
		Stdout: write,
		Stderr: write,
	}
	scanner := bufio.NewReader(read)
	buf := make([]byte, 0, 10*1024)
	op, err := connection.ExecContainer(instance, req, &args)
	if err != nil {
		return s, err
	}
	err = op.Wait()
	if err != nil {
		return s, err
	}
	s.ExitCode = int(op.Get().Metadata["return"].(float64))

	if s.ExitCode == 0 {
		go func() {
			time.Sleep(time.Millisecond * 1000)
			read.Close()
			write.Close()
		}()
		n, err := scanner.Read(buf[:cap(buf)])
		if err != nil {
			return s, err
		}
		s.Output = string(buf[:n])
	}
	return s, nil
}

func ClearInstance(instance string) error {
	req := api.ContainerExecPost{
		Command:     []string{"find", "/home/runner/", "-delete"},
		WaitForWS:   true,
		Interactive: false,
		User:        1500,
		Group:       1500,
		Cwd:         "/home/runner/",
	}
	args := lxd.ContainerExecArgs{
		Stdin:  nil,
		Stdout: nil,
		Stderr: nil,
	}
	op, err := connection.ExecContainer(instance, req, &args)
	if err != nil {
		return err
	}
	err = op.Wait()
	if err != nil {
		return err
	}
	return nil
}
