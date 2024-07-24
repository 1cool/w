/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"log"
)

// entityCmd represents the entity command
var entityCmd = &cobra.Command{
	Use:   "entity",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
