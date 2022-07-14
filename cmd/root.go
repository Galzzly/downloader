/*
Copyright © 2022 Liam Gallear

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	inFile      string
	concurrent  int
	destination string
)

var rootCmd = &cobra.Command{
	Use:   "downloader",
	Short: "Download from a list of URLs to desired location",
	Long: `The downloader tool will take in a list of URLs from a file
and download them to the specified location. 

Sub-commands will determine whether this is to store the files 
locally, or directly into HDFS.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Version = "0.1"
	rootCmd.PersistentFlags().StringVarP(&inFile, "file", "f", "", "file containing list of URLs")
	rootCmd.PersistentFlags().IntVarP(&concurrent, "concurrent", "c", 1, "number of concurrent downloads")
	rootCmd.PersistentFlags().StringVarP(&destination, "destination", "d", "", "desired destination to save the downloads")

	rootCmd.InitDefaultVersionFlag()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match
}
