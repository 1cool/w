/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	template2 "github.com/1cool/w/template"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"text/template"
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
		gen = &Generator{
			SnakeName: strcase.ToSnake(projectName),
			tmpl: template.Must(
				template.New("").
					Funcs(template.FuncMap{
						"ToSlash":     filepath.ToSlash,
						"FirstLetter": firstLetter,
					}).
					ParseFS(template2.TemplateDir, "tmpl/*.tmpl"),
			),
		}

		if fileIsExisted(gen.SnakeName) {
			return errors.New(fmt.Sprintf("the %s is existed not need init.\n", gen.SnakeName))
		}
		_, err := goModInit(gen.SnakeName)
		if err != nil {
			return errors.Wrapf(err, "go mod init error")
		}

		gen.WorkDir = filepath.Join("", gen.SnakeName)
		// variable
		gen.CamelName = "Example"
		gen.setVariable()

		if !fileIsExisted(gen.InternalDir) {
			err := os.MkdirAll(gen.InternalDir, os.ModePerm)
			if err != nil {
				return errors.Wrapf(err, "internal dir init error")
			}
		}

		gen.InjectInterfaceImpl = "{{ .InjectInterfaceImpl }}"
		gen.InjectInterface = "{{ .InjectInterface }}"

		_, err = gen.NewEntSchema()
		if err != nil {
			return errors.Wrapf(err, "ent new error")
		}

		_, err = gen.EntGenerate()
		if err != nil {
			return errors.Wrapf(err, "go generate error")
		}

		err = gen.InitGenerate("repository.go", gen.RepositoryDir, "repository")
		if err != nil {
			return errors.Wrapf(err, "repository init error")
		}

		err = gen.InitGenerate("service.go", gen.ServiceDir, "service")
		if err != nil {
			return errors.Wrapf(err, "service init error")
		}

		err = gen.InitGenerate("gintransport.go", gen.HandlerDir, "handler.tmpl")
		if err != nil {
			return errors.Wrapf(err, "handler init error")
		}

		err = gen.InitGenerate("example.go", gen.HandlerDir, "handler_example.tmpl")
		if err != nil {
			return errors.Wrapf(err, "handler example init error")
		}

		err = gen.InitGenerate("config.go", gen.ModelDir, "config.tmpl")
		if err != nil {
			return errors.Wrapf(err, "config init error")
		}

		err = gen.InitGenerate("viper.go", gen.ConfigDir, "viper.tmpl")
		if err != nil {
			return errors.Wrapf(err, "viper init error")
		}

		err = gen.InitGenerate("database.go", gen.DatabaseDir, "database.tmpl")
		if err != nil {
			return errors.Wrapf(err, "database init error")
		}

		err = gen.InitGenerate("mysql.go", gen.DatabaseDir, "mysql.tmpl")
		if err != nil {
			return errors.Wrapf(err, "mysql init error")
		}

		err = gen.InitGenerate("main.go", gen.WorkDir, "main.tmpl")
		if err != nil {
			return errors.Wrapf(err, "main init error")
		}

		err = gen.InitGenerate("config.yaml", gen.WorkDir, "configyaml.tmpl")
		if err != nil {
			return errors.Wrapf(err, "config.yaml init error")
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
