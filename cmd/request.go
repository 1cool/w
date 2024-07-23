/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

// requestCmd represents the request command
var requestCmd = &cobra.Command{
	Use:   "request",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("project_name needs to be provided")
		}
		gen := NewGenerate(args[0])
		moduleName, err := ReadModuleNameFromGoModFile()
		if err != nil {
			return err
		}
		gen.Module = moduleName
		gen.setDir()

		gen.Generate(filepath.Join(gen.RequestDir, "", ActionCreate))
		fmt.Println("request called")
	},
}

func init() {
	newCmd.AddCommand(requestCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// requestCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// requestCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func (gen *Generator) newRequest() error {
	gen.setDir()
	gen.tmpl.ExecuteTemplate(&gen.TemplateInterface, "repository_interface", gen)

}
