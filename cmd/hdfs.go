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
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Galzzly/downloader/internal"
	"github.com/colinmarc/hdfs/v2"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
)

// hdfsCmd represents the hdfs command
var hdfsCmd = &cobra.Command{
	Use:   "hdfs",
	Short: "Download files to a HDFS path",
	Long: `Take in the list of URLs to download and
save them to a HDFS path.`,
	Run: func(cmd *cobra.Command, args []string) {
		runHDFS()
	},
}

func init() {
	rootCmd.AddCommand(hdfsCmd)
}

func runHDFS() {
	// Get the list of addresses to cycle through
	addressList, err := internal.GetAddresses(inFile)
	if err != nil {
		fmt.Fprint(os.Stderr, "Error getting addresses: ", err)
		return
	}

	// Get the HDFS connection
	hdfscli, err := internal.ConnectToNamenode()
	if err != nil {
		fmt.Fprint(os.Stderr, "Error connecting to HDFS: ", err)
		return
	}

	// Check if the HDFS destination directory exists, create if not
	if _, err := hdfscli.Stat(destination); err != nil {
		if err := hdfscli.MkdirAll(destination, 0755); err != nil {
			fmt.Fprint(os.Stderr, "Error creating destination directory: ", err)
		}
	}

	// Set up our progress bars
	p := mpb.New(mpb.PopCompletedMode())
	bar := p.Add(
		int64(len(addressList)),
		mpb.NewBarFiller(
			mpb.BarStyle().Lbound("╢").Filler("▌").Tip("▌").Padding("░").Rbound("╟"),
		),
		mpb.BarNoPop(),
		mpb.BarFillerClearOnComplete(),
		mpb.PrependDecorators(
			decor.Name("Downloading", decor.WC{W: 11}),
		),
		mpb.AppendDecorators(
			decor.OnComplete(decor.CountersNoUnit("%d /%d", decor.WC{W: 10, C: decor.DidentRight}), "Done!"),
			decor.Elapsed(decor.ET_STYLE_GO, decor.WCSyncSpace),
		),
	)

	// Set our channels up per concurrency value
	workers := make(chan int, concurrent)

	// Do the work to download the files
	for _, i := range addressList {
		go downloadHdfsFile(bar, workers, hdfscli, i, destination)
	}

	bar.Wait()
	p.Wait()
	fmt.Println(internal.Complete())
}

func downloadHdfsFile(bar *mpb.Bar, worker chan int, hdfscli *hdfs.Client, address, destination string) {
	worker <- 1
	filename := address[strings.LastIndex(address, "/")+1:]
	target := filepath.Join(destination, filename)

	out, err := hdfscli.Create(target)
	if err != nil {
		bar.Increment()
		failures++
		<-worker
		return
	}

	resp, err := internal.DownloadFile(address)
	if err != nil {
		bar.Increment()
		failures++
		<-worker
		return
	}
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		out.Close()
		bar.Increment()
		failures++
		<-worker
		return
	}

	out.Close()
	bar.Increment()
	<-worker
}
