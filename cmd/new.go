/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	template2 "github.com/1cool/wo/template"
	"github.com/spf13/cobra"
	"path/filepath"
	"text/template"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "new command for exit init project.",
	Long: `For exit project run cli new command generate code. For example:

w new [command] [name]`,
	RunE: func(cmd *cobra.Command, args []string) error {
		gen = &Generator{
			tmpl: template.Must(
				template.New("new").
					Funcs(template.FuncMap{
						"ToSlash": filepath.ToSlash,
					}).
					ParseFS(template2.TemplateDir, "tmpl/*.tmpl")),
		}
		fmt.Println("new called")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
