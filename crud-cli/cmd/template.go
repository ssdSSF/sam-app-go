/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

var sets []string
var file string

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Templating the sam template.yaml",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		t, err := template.ParseFiles(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "err parse '%s': %s\n", file, err)
			os.Exit(1)
		}
		keyValues := make(map[string]string)
		for _, set := range sets {
			keyValues[strings.Split(set, "=")[0]] = strings.Split(set, "=")[1]
		}
		t.Execute(os.Stdout, keyValues)
	},
}

func init() {
	rootCmd.AddCommand(templateCmd)

	templateCmd.Flags().StringArrayVar(&sets, "set", sets, "key values for the sam environment variables")
	templateCmd.Flags().StringVarP(&file, "file", "f", file, "file path of the template.goyaml")
}
