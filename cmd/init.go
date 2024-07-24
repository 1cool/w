/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

// initCmd represents the init command
// 1、mod init
// 2、ent init
// 3、repository init
// 4、service init
// 5、handler init
var initCmd = &cobra.Command{
	Use:   "init [project_name]",
	Short: "init one golang project",
	Long: `init one golang project. For example:

w init project`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("project_name needs to be provided")
		}

		projectName := args[0]
		gen = NewGenerate(projectName)

		if fileIsExisted(gen.SnakeName) {
			return errors.New(fmt.Sprintf("the %s is existed not need init.\n", gen.SnakeName))
		}
		_, err := goModInit(gen.SnakeName)
		if err != nil {
			return errors.Wrapf(err, "go mod init error")
		}

		gen.WorkDir = filepath.Join("", gen.SnakeName)
		gen.Module = gen.SnakeName
		// variable
		gen.setDir()
		_, err = gen.NewEntSchema()
		if err != nil {
			return errors.Wrapf(err, "ent new error")
		}

		_, err = gen.EntGenerate()
		if err != nil {
			return errors.Wrapf(err, "go generate error")
		}

		err = gen.Generate(filepath.Join(gen.RepositoryDir, "repository.go"), "repository", ActionCreate)
		if err != nil {
			return errors.Wrapf(err, "repository init error")
		}

		err = gen.Generate(filepath.Join(gen.ServiceDir, "service.go"), "service", ActionCreate)
		if err != nil {
			return errors.Wrapf(err, "service init error")
		}

		err = gen.Generate(filepath.Join(gen.HandlerDir, "gintransport.go"), "handler.tmpl", ActionCreate)
		if err != nil {
			return errors.Wrapf(err, "handler init error")
		}

		err = gen.Generate(filepath.Join(gen.HandlerDir, "example.go"), "handler_example.tmpl", ActionCreate)
		if err != nil {
			return errors.Wrapf(err, "config init error")
		}

		err = gen.Generate(filepath.Join(gen.ModelDir, "config.go"), "config.tmpl", ActionCreate)
		if err != nil {
			return errors.Wrapf(err, "config init error")
		}

		err = gen.Generate(filepath.Join(gen.ModelDir, "pagination.go"), "pagination.tmpl", ActionCreate)
		if err != nil {
			return errors.Wrapf(err, "pagination init error")
		}

		err = gen.Generate(filepath.Join(gen.ConfigDir, "viper.go"), "viper.tmpl", ActionCreate)
		if err != nil {
			return errors.Wrapf(err, "viper init error")
		}

		err = gen.Generate(filepath.Join(gen.DatabaseDir, "database.go"), "database.tmpl", ActionCreate)
		if err != nil {
			return errors.Wrapf(err, "database init error")
		}

		err = gen.Generate(filepath.Join(gen.DatabaseDir, "mysql.go"), "mysql.tmpl", ActionCreate)
		if err != nil {
			return errors.Wrapf(err, "mysql init error")
		}

		err = gen.Generate(filepath.Join(gen.WorkDir, "main.go"), "main.tmpl", ActionCreate)
		if err != nil {
			return errors.Wrapf(err, "main init error")
		}

		err = gen.Generate(filepath.Join(gen.WorkDir, "config.yaml"), "configyaml.tmpl", ActionCreate)
		if err != nil {
			return errors.Wrapf(err, "config.yaml init error")
		}

		err = gen.Generate(filepath.Join(gen.InternalDir, "error.go"), "error.tmpl", ActionCreate)
		if err != nil {
			return errors.Wrapf(err, "error.yaml init error")
		}

		err = gen.Generate(filepath.Join(gen.InternalDir, "constant.go"), "constant.tmpl", ActionCreate)
		if err != nil {
			return errors.Wrapf(err, "constant.yaml init error")
		}

		fmt.Println("init successful")
		fmt.Println("please cd your project dirname,then edit your config file `config.yaml`")
		fmt.Println("finally please run `go mod tidy && go run main.go`")
		return nil
	},
}

func goModInit(moduleName string) (string, error) {
	dir, _ := os.Getwd()
	projectDir := filepath.Join(dir, moduleName)

	// 项目目录
	if err := os.MkdirAll(projectDir, os.ModePerm); err != nil {
		return "", errors.Wrapf(err, "mkdir %s failed", moduleName)
	}

	out, err := runCmdCommand("go", projectDir, []string{
		"mod",
		"init",
		moduleName,
	}...)

	if err != nil {
		return "", err
	}

	return out, nil
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
