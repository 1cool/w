/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"path/filepath"
)

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
	Use:   "service [service_name]",
	Short: "new service",
	Long: `new service. For example:

gapi new service user.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			log.Fatalln("service name needs to be provided")
		}

		gen := NewGenerate(args[0])
		moduleName, err := ReadModuleNameFromGoModFile()
		if err != nil {
			return err
		}
		gen.Module = moduleName
		gen.setDir()

		err = gen.NewService()
		if err != nil {
			return err
		}

		fmt.Println("new service successful", gen.SnakeName)
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

func (gen *Generator) NewService() error {
	err := gen.Generate(filepath.Join(gen.ServiceDir, gen.SnakeName+".go"), "service_entity", ActionCreate)
	if err != nil {
		return err
	}

	gen.UpdatePosition = "// InjectInterface"
	err = gen.Generate(filepath.Join(gen.ServiceDir, "service.go"), "service_interface", ActionUpdate)
	if err != nil {
		return err
	}

	gen.UpdatePosition = "// InjectInterfaceImpl"
	return gen.Generate(filepath.Join(gen.ServiceDir, "service.go"), "service_interface_impl", ActionUpdate)
}
