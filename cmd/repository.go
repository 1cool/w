/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	template2 "github.com/1cool/w/template"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"path/filepath"
	"text/template"
)

// repositoryCmd represents the repository command
var repositoryCmd = &cobra.Command{
	Use:   "repository [repository_name] [func_name]",
	Short: "new repository",
	Long: `new repository. For example:

new repository user addUser`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("repository_name and func_name needs to be provided")
		}

		name := args[0]
		gen = &Generator{
			tmpl: template.Must(
				template.New("new").
					Funcs(template.FuncMap{
						"ToSlash":     filepath.ToSlash,
						"FirstLetter": firstLetter,
					}).
					ParseFS(template2.TemplateDir, "tmpl/*.tmpl")),
		}

		moduleName, err := ReadModuleNameFromGoModFile()
		if err != nil {
			return err
		}

		gen.Module = moduleName
		gen.CamelName = strcase.ToCamel(name)
		gen.LowerCamelName = strcase.ToLowerCamel(name)
		gen.SnakeName = strcase.ToSnake(name)
		gen.FuncName = strcase.ToCamel(args[1])

		gen.setVariable()
		gen.Type = TypeRepository

		err = gen.NewRepository()
		if err != nil {
			return err
		}

		fmt.Println("new repository successful", gen.SnakeName, gen.FuncName)
		return nil
	},
}

func init() {
	newCmd.AddCommand(repositoryCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// repositoryCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// repositoryCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// NewRepository
// 1. generate CamelNameRepository file
// 2. generate NewCamelNameRepository interface in Repository Interface file
// 3. generate NewCamelNameRepositoryImpl function in Repository Interface file
func (gen *Generator) NewRepository() error {
	gen.TargetFuncFile = filepath.Join(gen.RepositoryDir, gen.SnakeName+".go")
	gen.TargetInterfaceFile = filepath.Join(gen.RepositoryDir, "repository.go")
	gen.InjectInterfaceEntity = bytes.Buffer{}
	gen.TemplateFunc = bytes.Buffer{}
	gen.TemplateInterface = bytes.Buffer{}
	gen.InjectInterfaceImpl = "{{ .InjectInterfaceImpl }}"
	gen.InjectInterface = "{{ .InjectInterface }}"
	gen.InjectHereImpl = "{{ .InjectHereImpl }}"

	fmt.Println(fmt.Sprintf("%+v", gen))
	//err := gen.tmpl.ExecuteTemplate(&gen.TemplateInterface, "repository_interface.tmpl", gen)
	//if err != nil {
	//	return err
	//}

	err := gen.tmpl.ExecuteTemplate(&gen.InjectInterfaceEntity, "repository_entity", gen)
	if err != nil {
		return err
	}

	err = WriteToFile(gen.TargetFuncFile, gen.InjectInterfaceEntity.Bytes())
	if err != nil {
		return err
	}

	return gen.NewGenerate()
}
