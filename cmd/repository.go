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

// repositoryCmd represents the repository command
var repositoryCmd = &cobra.Command{
	Use:   "repository [repository_name]",
	Short: "new repository",
	Long: `new repository. For example:

new repository user addUser`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			log.Fatalln("repository name needs to be provided")
		}

		gen := NewGenerate(args[0])
		moduleName, err := ReadModuleNameFromGoModFile()
		if err != nil {
			return err
		}
		gen.Module = moduleName
		gen.setDir()

		err = gen.NewRepository()
		if err != nil {
			return err
		}
		fmt.Println("new repository successful", gen.SnakeName)
		return nil
	},
}

func init() {
	//newCmd.AddCommand(repositoryCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// repositoryCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// repositoryCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func (gen *Generator) NewRepository() error {
	err := gen.Generate(filepath.Join(gen.ModelDir, gen.SnakeName+".go"), "model", ActionCreate)
	if err != nil {
		return err
	}

	err = gen.Generate(filepath.Join(gen.RepositoryDir, gen.SnakeName+".go"), "repository_entity", ActionCreate)
	if err != nil {
		return err
	}

	gen.UpdatePosition = "// InjectInterface"
	err = gen.Generate(filepath.Join(gen.RepositoryDir, "repository.go"), "repository_interface", ActionUpdate)
	if err != nil {
		return err
	}

	gen.UpdatePosition = "// InjectInterfaceImpl"
	return gen.Generate(filepath.Join(gen.RepositoryDir, "repository.go"), "repository_interface_impl", ActionUpdate)
}
