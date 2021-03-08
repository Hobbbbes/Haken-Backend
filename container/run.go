package container

import (
	"bufio"
	"fmt"
	"io"
	"os"

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

type Status struct {
	ExitCode int
	Output   string
}

//PrepareExecution prepares the execution enviroment by copying source files and executing the PreLaunchTask
func PrepareExecution(SourcePath string, lang datastructures.Language, instance string) (Status, error) {
	var s Status
	f, err := os.Open(SourcePath)
	if err != nil {
		return s, err
	}
	args := lxd.InstanceFileArgs{
		Content: f,
		UID:     1500,
		GID:     1500,
		Mode:    644,
	}
	err = connection.CreateInstanceFile(instance, "/home/runner/main"+lang.Abbreviation, args)
	if err != nil {
		return s, err
	}
	if lang.PreLaunchTask == "" {
		s.ExitCode = -1
		s.Output = ""
		return s, nil
	}
	req := api.ContainerExecPost{
		Command:     []string{fmt.Sprintf(lang.PreLaunchTask, "main"+lang.Abbreviation)},
		WaitForWS:   true,
		Interactive: false,
		User:        1500,
		Group:       1500,
		Cwd:         "/home/runner/",
	}
	read, write := io.Pipe()
	args2 := lxd.ContainerExecArgs{
		Stdin:  nil,
		Stdout: write,
		Stderr: write,
	}
	scanner := bufio.NewScanner(read)
	op, err := connection.ExecContainer("h2", req, &args2)
	if err != nil {
		return s, err
	}
	// Wait for it to complete
	err = op.Wait()
	if err != nil {
		return s, err
	}
	s.ExitCode = op.Get().Metadata["return"].(int)
	if s.ExitCode == 0 {
		return s, nil
	}
	s.Output = ""
	for scanner.Scan() {
		s.Output += scanner.Text()
		s.Output += "\n"
	}
	return s, nil
}
