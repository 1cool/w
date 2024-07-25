package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"strings"
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

func firstLetter(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToLower(string(s[0]))
}

func WriteToFile(newFile string, fileData []byte) error {
	if err := os.WriteFile(newFile, fileData, os.ModePerm); err != nil {
		return fmt.Errorf("writing file %s: %w", newFile, err)
	}

	return nil
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

func ReadModuleNameFromGoModFile() (string, error) {
	if !fileIsExisted("go.mod") {
		return "", errors.New("请在go.mod同级目录下执行命令")
	}

	file, err := os.Open("go.mod")
	if err != nil {
		return "", errors.Wrapf(err, "open go.mod file error")
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "module ") {
			modulePath := strings.TrimSpace(strings.TrimPrefix(line, "module "))
			return modulePath, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", errors.Wrapf(err, "Error reading go.mod file error")
	}

	return "", errors.New("go.mod file not hava module")
}
