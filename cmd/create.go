package cmd

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type command struct {
	Class string
	Description string
	Command string
	Flags []flag
}

type flag struct {
	Flag string
	Description string
	Value string
}

type foo struct {
	AppName string
	Package string
	Commands []command
}

func init() {
	rootCmd.AddCommand(createCmd)

}

var createCmd = &cobra.Command{
	Use:   "create [name] [dir]",
	Short: "Create a new severell project",
	Long:  `All software has versions. This is Hugo's`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Base Package (Group ID): ")
		basePackage, _ := reader.ReadString('\n')
		basePackage = strings.Replace(basePackage, string('\n'), "", 1)

		createDir(args[0], basePackage)
		downloadZip(args[0])
		dest := strings.Replace(args[0], ".", "", 1)
		fmt.Println("Dest:" + dest)
		_, err := unzip(args[0] + "/test.zip", dest, basePackage, args[0])
		fmt.Println(err)

		data := foo {
			Package: basePackage,
		}

		file, _ := json.MarshalIndent(data, "", " ")

		_ = ioutil.WriteFile("severell.json", file, 0644)

	},
}

func createDir(dir string, basePackage string) {
	if dir != "." {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.Mkdir(dir, os.ModePerm)
		}
	}
}

func downloadZip(dir string) {
	specUrl := "https://github.com/mitchdennett/severell-framework/archive/master.zip"
	resp, err := http.Get(specUrl)
	if err != nil {
		fmt.Printf("err: %s", err)
	}


	defer resp.Body.Close()
	fmt.Println("status", resp.Status)
	if resp.StatusCode != 200 {
		return
	}

	// Create the file
	out, err := os.Create(dir+ "/test.zip")
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
}

func unzip(src string, dest string, basePackage string, name string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on

		fpath := filepath.Join(dest, strings.TrimPrefix(f.Name, "severell-framework-master/"))
		fmt.Println(fpath)
		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		//if fpath != dest && !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
		//	return filenames, fmt.Errorf("%s: illegal file path", fpath)
		//}
		replaceDest := dest
		if dest != "" {
			replaceDest = dest + "/"
		}

		if strings.HasPrefix(fpath, replaceDest + "src/main/java") ||  strings.HasPrefix(fpath, replaceDest + "src/test/java") {
			fpath = strings.Replace(fpath, replaceDest + "src/main/java", replaceDest + "src/main/java/" + strings.ReplaceAll(basePackage, ".", "/"), 1)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		if !strings.Contains(filepath.Base(fpath), ".mustache") {
			buf := new(bytes.Buffer)
			buf.ReadFrom(rc)

			tmpl, err := template.New("test").Parse(buf.String())
			if err != nil {
				fmt.Println(err)
			}
			err = tmpl.Execute(outFile, &foo{Package: strings.ReplaceAll(basePackage, string('\n'), ""), AppName: name})
			if err != nil {
				fmt.Println(err)
			}
		} else {
			_, err = io.Copy(outFile, rc)
		}

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}