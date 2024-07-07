package cmd

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"os/exec"
)

func fileIsExisted(directory string) bool {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		return false
	}
	return true
}

func entNew(schema string) (string, error) {
	return runCmdCommand("go", "", []string{
		"run",
		"-mod=mod",
		"entgo.io/ent/cmd/ent",
		"new",
		schema,
	}...)
}

func WriteToFile(newFile string, fileData []byte) error {
	if err := os.WriteFile(newFile, fileData, os.ModePerm); err != nil {
		return fmt.Errorf("writing file %s: %w", newFile, err)
	}

	return nil
}

func generate() (string, error) {
	return runCmdCommand("go", "", []string{
		"generate",
		"./ent",
	}...)
}

func runCmdCommand(command string, dir string, args ...string) (string, error) {
	var (
		o bytes.Buffer
		e bytes.Buffer
	)

	cmd := exec.Command(command, args...)
	cmd.Stdout = &o
	cmd.Stderr = &e

	if dir != "" {
		cmd.Dir = dir
	}

	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", errors.WithMessagef(err, "命令执行失败;exitcode:%d;message:%s", exitErr.ExitCode(), string(exitErr.Stderr)+":"+e.String())
		} else {
			return "", errors.WithMessage(err, "执行命令错误")
		}
	}
	return o.String(), nil
}
