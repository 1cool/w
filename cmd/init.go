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
	"log"
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
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatalln("project_name needs to be provided")
		}

		projectName := args[0]
		gen = &Generator{
			tmpl: template.Must(
				template.New("").
					Funcs(template.FuncMap{
						"ToSlash": filepath.ToSlash,
					}).
					ParseFS(template2.TemplateDir, "tmpl/*.tmpl"),
			),
		}
		gen.SnakeName = strcase.ToSnake(projectName)
		if fileIsExisted(gen.SnakeName) {
			log.Fatalf("the %s is existed not need init.\n", gen.SnakeName)
		}
		_, err := goModInit(gen.SnakeName)
		if err != nil {
			log.Fatalln("go mod init error:", err)
		}

		gen.WorkDir = filepath.Join("", gen.SnakeName)
		internal := filepath.Join(gen.WorkDir, "internal")
		if !fileIsExisted(internal) {
			err := os.MkdirAll(internal, os.ModePerm)
			if err != nil {
				log.Fatalln("internal dir init error:", err)
			}
		}

		gen.InternalDir = internal
		gen.ServiceDir = filepath.Join(gen.InternalDir, "service")
		gen.RepositoryDir = filepath.Join(gen.InternalDir, "repository")
		gen.ModelDir = filepath.Join(gen.InternalDir, "model")
		gen.ConfigDir = filepath.Join(gen.InternalDir, "config")
		gen.DatabaseDir = filepath.Join(gen.InternalDir, "database")
		gen.HandlerDir = filepath.Join(gen.InternalDir, "httptransport")
		gen.InjectInterfaceImpl = "{{ .InjectInterfaceImpl }}"
		gen.InjectInterface = "{{ .InjectInterface }}"

		gen.CamelName = "User"

		_, err = gen.NewEntSchema()
		if err != nil {
			log.Fatalln("ent new error:", err)
		}

		_, err = gen.EntGenerate()
		if err != nil {
			log.Fatalln("go generate error:", err)
		}

		err = gen.InitGenerate("repository.go", gen.RepositoryDir, "repository_interface.tmpl")
		if err != nil {
			log.Fatalln("repository init error:", err)
		}

		err = gen.InitGenerate("service.go", gen.ServiceDir, "service_interface.tmpl")
		if err != nil {
			log.Fatalln("service init error:", err)
		}

		err = gen.InitGenerate("gintransport.go", gen.HandlerDir, "handler.tmpl")
		if err != nil {
			log.Fatalln("handler init error:", err)
		}

		err = gen.InitGenerate("example.go", gen.HandlerDir, "handler_example.tmpl")
		if err != nil {
			log.Fatalln("handler example init error:", err)
		}

		err = gen.InitGenerate("config.go", gen.ModelDir, "config.tmpl")
		if err != nil {
			log.Fatalln("model init error:", err)
		}

		err = gen.InitGenerate("viper.go", gen.ConfigDir, "viper.tmpl")
		if err != nil {
			log.Fatalln("config init error:", err)
		}

		err = gen.InitGenerate("database.go", gen.DatabaseDir, "database.tmpl")
		if err != nil {
			log.Fatalln("database init error:", err)
		}

		err = gen.InitGenerate("mysql.go", gen.DatabaseDir, "mysql.tmpl")
		if err != nil {
			log.Fatalln("mysql init error:", err)
		}

		err = gen.InitGenerate("main.go", gen.WorkDir, "main.tmpl")
		if err != nil {
			log.Fatalln("main init error:", err)
		}

		err = gen.InitGenerate("config.yaml", gen.WorkDir, "configyaml.tmpl")
		if err != nil {
			log.Fatalln("main init error:", err)
		}

		fmt.Println("init successful")
		fmt.Println("please cd your project dirname,then edit your config file `config.yaml`")
		fmt.Println("finally please run `go mod tidy && go run main.go`")
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
