/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"log"
	"path/filepath"
)

// entityCmd represents the entity command
var entityCmd = &cobra.Command{
	Use:   "entity",
	Short: "new entity [name]",
	Long:  `generate model ent schema repository service handler request response`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			log.Fatalln("entity name needs to be provided")
		}

		gen := NewGenerate(args[0])
		moduleName, err := ReadModuleNameFromGoModFile()
		if err != nil {
			return err
		}
		gen.Module = moduleName
		gen.setDir()

		_, err = gen.NewEntSchema()
		if err != nil {
			return errors.Wrapf(err, "ent new error")
		}

		_, err = gen.EntGenerate()
		if err != nil {
			return errors.Wrapf(err, "go generate error")
		}

		err = gen.NewRepository()
		if err != nil {
			return err
		}

		err = gen.NewService()
		if err != nil {
			return err
		}

		err = gen.NewHandler()
		if err != nil {
			return err
		}

		err = gen.NewRequest()
		if err != nil {
			return err
		}
		err = gen.NewResponse()
		if err != nil {
			return err
		}

		err = gen.NewRouter()
		if err != nil {
			return err
		}

		fmt.Println("new entity successful", gen.SnakeName)
		return nil
	},
}

func init() {
	newCmd.AddCommand(entityCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// entityCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// entityCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func (gen *Generator) NewRouter() error {
	gen.UpdatePosition = "// InjectRouter"
	return gen.Generate(filepath.Join(gen.HandlerDir, "gintransport.go"), "router.tmpl", ActionUpdate)
}
