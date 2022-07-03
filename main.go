package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	. "internal/du"
	. "internal/tui"
)

/*
	This var sets up the root command and then all other commands. The root command, according to Cobra's structure, is the first thing we hit when we run the program.
	Imagine it as an automatic constructor that's allowing us to run an instance of this program.
*/
var rootCmd = &cobra.Command{
	Use: "put usage example here",
	//TraverseChildren: true,
	Short: "This program shows disk usage",
	Long:  "Put longer version of Short here",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(args)
		v, _ := cmd.Flags().GetBool("version")
		if v {
			version()
		} else {
			//log.Fatal("version unavailable")
			fmt.Println("version not checked")
		}
		l := logFlag
		if l != "" {
			logFile, err = os.OpenFile(l, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
			if err != nil {
				err = fmt.Errorf("error opening log file: %w", err)
			}
		}
		defer func() {
			cerr := logFile.Close()
			if err == nil {
				err = cerr
			}
		}()
		log.SetOutput(logFile)
		o := outputFlag
		if o == "-" {
			outputFile = os.Stdout
		} else {
			outputFile, err = os.OpenFile(o, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
			if err != nil {
				err = fmt.Errorf("error setting output file: %w", err)
			}
		}
		if len(args) == 1 {
			dir, _ = filepath.Abs(args[0])
		} else {
			dir = "."
		}
		f := inputFile
		if f != "" {
			//needs to have logic for setting input file
			fmt.Println(f)
		}
		e := extendedFlag
		if e {
			fmt.Println("Needs implementation of extended information output")
		}
		ic := icFlag
		if ic {
			fmt.Println("Needs implementation of ignore configuration")
		}

	},
}

var (
	icFlag       bool
	extendedFlag bool
	versionFlag  bool
	inputFile    string
	outputFlag   string
	outputFile   io.Writer
	logFlag      string
	logFile      *os.File
	err          error
	dir          string
)

func version() {
	fmt.Println("Version goes here")
}

func init() {
	flags := rootCmd.Flags()

	flags.StringVarP(&outputFlag, "output-file", "o", "", "-o [FILE] defines file for data output")
	flags.StringVarP(&inputFile, "input-file", "f", "", "-f [FILE] defines file for data input")
	flags.BoolVarP(&versionFlag, "version", "v", false, "-v shows the current version of godu")
	flags.BoolVarP(&extendedFlag, "extended", "e", false, "-e enables extended information mode")
	flags.BoolVar(&icFlag, "ignore-config", false, "--ignore-config prevents godu from attempting to load any configuration files")

}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

	files, err := ListFilesRecursivelyInParallel(".")
	if err != nil {
		log.Fatalln(err)
	}

	if len(files) > 0 {
		log.Println(files[0])
	}

	p := tea.NewProgram(InitialModel(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}

func ListFilesRecursivelyInParallel(s string) {
	panic("unimplemented")
}
