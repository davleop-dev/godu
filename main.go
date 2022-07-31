package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	du "internal/du"
	tui "internal/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

const godu_version = "v0.1.0a"

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

		o := outputFlag
		if o == "-" {
			outputFile = os.Stdout
		} else {
			outputFile, err = os.OpenFile(o, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
			if err != nil {
				err = fmt.Errorf("error setting output file: %w", err)
			}
		}

		defer func() {
			cerr := logFile.Close()
			if err == nil {
				err = cerr
			}
		}()
		log.SetOutput(logFile)

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
		x := xFlag
		if x {
			fmt.Println("Needs implementation of filesystem boundaries")
		}
		c := cfsFlag
		if c {
			fmt.Println("Needs implementation of filesystem boundaries")
		}
		ex := exclude
		if len(ex) != 0 {
			for i := 0; i < len(ex); i++ {
				fmt.Println(ex[i])
			}
		}
		bX := bigXFlag
		if bX != "" {
			fmt.Println("Needs implementation of Exclusion Files")
		}
		sym := symLinkFlag
		if sym {
			fmt.Println("Needs implementation of symLink following")
		}
	},
}

var (
	symLinkFlag  bool
	bigXFlag     string
	exclude      []string
	cfsFlag      bool
	xFlag        bool
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
	fmt.Println(godu_version)
}

func init() {
	flags := rootCmd.Flags()

	flags.StringVarP(&outputFlag, "output-file", "o", "", "-o [FILE] defines file for data output")
	flags.StringVarP(&inputFile, "input-file", "f", "", "-f [FILE] defines file for data input")
	flags.BoolVarP(&versionFlag, "version", "v", false, "-v shows the current version of godu")
	flags.BoolVarP(&extendedFlag, "extended", "e", false, "-e enables extended information mode")
	flags.BoolVar(&icFlag, "ignore-config", false, "--ignore-config prevents godu from attempting to load any configuration files")
	flags.BoolVarP(&xFlag, "one-file-system", "x", false, "-x prevents godu from crossing filesystem boundaries, i.e. only count files and directories on the same filesystem as the directory being scanned")
	flags.BoolVar(&cfsFlag, "cross-file-system", false, "--cross-file-system allows godu to cross filesystem boundaries. This is the default, but can be specified to overrule a previously given '-x'")
	flags.StringArrayVar(&exclude, "exclude", exclude, "--exclude [PATTERN] excludes files that match PATTERN. The files will still be displayed by default, but are not counted towards the disk usage statistics. This argument can be added multiple times to add more patterns.")
	flags.StringVarP(&bigXFlag, "exclude-from", "X", "", "-X [FILE], --exclude-from [FILE] Exclude files that match any pattern in FILE. Patterns should be separated by a newline.")
	flags.BoolVarP(&symLinkFlag, "follow-symlinks", "L", false, "-L, --follow-symlinks follows symlinks and counts the size of the file they point to.")

}

func crappyCalculation(files []du.File) (map[string]int64, int64) {
	drsz := make(map[string]int64)
	totalSz := int64(0)

	for _, file := range files {
		if _, ok := drsz[file.HighDir]; ok {
			drsz[file.HighDir] += file.Size
		} else {
			drsz[file.HighDir] = file.Size
		}
		totalSz += file.Size
	}

	return drsz, totalSz
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

	hidden := true
	defaultOrdering := tui.Size
	directoryFirst := false
	desc := true

	files, err := du.ListFilesRecursivelyInParallel(dir)
	if err != nil {
		log.Fatalln(err)
	}

	drsz, totsz := crappyCalculation(files)

	initialModel := tui.Model{
		CurrentDirectory: dir,
		ShowHidden:       hidden,
		ListOrder:        defaultOrdering,
		Descending:       desc,
		DirectoryFirst:   directoryFirst,
		Files:            files,
		DirSz:            drsz,
		TotalSz:          totsz,
		Version:          godu_version,
	}

	p := tea.NewProgram(tui.NewModel(initialModel), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
