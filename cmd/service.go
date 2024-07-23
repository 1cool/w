/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	template2 "github.com/1cool/w/template"
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
	"log"
	"path/filepath"
	"text/template"
)

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
	Use:   "service [service_name] [func_name]",
	Short: "new service",
	Long: `new service. For example:

gapi new service user.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			log.Fatalln("service_name and func_name needs to be provided")
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
		gen.setDir()
		gen.SnakeName = strcase.ToSnake(name)

		//err = gen.NewService()

		if err != nil {
			return err
		}

		//fmt.Println("new service successful", gen.SnakeName, gen.FuncName)
		return nil
	},
}

func init() {
	newCmd.AddCommand(serviceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serviceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serviceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

//func (gen *Generator) NewService() error {
//	gen.TargetFuncFile = filepath.Join(gen.ServiceDir, gen.SnakeName+".go")
//	gen.TargetInterfaceFile = filepath.Join(gen.ServiceDir, "service.go")
//	gen.InjectInterfaceEntity = bytes.Buffer{}
//	gen.TemplateFunc = bytes.Buffer{}
//	gen.TemplateInterface = bytes.Buffer{}
//	gen.InjectInterfaceImpl = "{{ .InjectInterfaceImpl }}"
//	gen.InjectInterface = "{{ .InjectInterface }}"
//
//	err := gen.tmpl.ExecuteTemplate(&gen.InjectInterfaceEntity, "service_entity", gen)
//	if err != nil {
//		return err
//	}
//	err = WriteToFile(gen.TargetFuncFile, gen.InjectInterfaceEntity.Bytes())
//	if err != nil {
//		return err
//	}
//
//	interfaceData, err := os.ReadFile(gen.TargetInterfaceFile)
//	if err != nil {
//		return err
//	}
//	err = gen.tmpl.ExecuteTemplate(&gen.TemplateInterface, "service_interface", gen)
//	if err != nil {
//		return err
//	}
//	newInterfaceData := strings.Replace(string(interfaceData), "// {{ .InjectInterface }}", gen.TemplateInterface.String(), 1)
//	if err != nil {
//		return err
//	}
//	err = gen.tmpl.ExecuteTemplate(&gen.TemplateFunc, "service_interface_impl", gen)
//	if err != nil {
//		return err
//	}
//	newInterfaceData = strings.Replace(newInterfaceData, "// {{ .InjectInterfaceImpl }}", gen.TemplateFunc.String(), 1)
//	return WriteToFile(gen.TargetInterfaceFile, []byte(newInterfaceData))
//}
