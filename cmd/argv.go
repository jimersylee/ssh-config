package cmd

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type Args struct {
	ToYAML   bool
	ToSSH    bool
	ToJSON   bool
	Src      string
	Dest     string
	ShowHelp bool
}

func ParseArgs() (result Args) {
	toYAML := flag.Bool("to-yaml", false, "Convert SSH config(Text/JSON) to YAML")
	toSSH := flag.Bool("to-ssh", false, "Convert SSH config(YAML/JSON) to YAML")
	toJSON := flag.Bool("to-json", false, "Convert SSH config(YAML/Text) to JSON")
	src := flag.String("src", "", "Source file or directories path, valid when using non-pipeline mode")
	dest := flag.String("dest", "", "Destination file path, valid when using non-pipeline mode")
	showHelp := flag.Bool("help", false, "Show help")

	flag.Parse()

	return Args{
		ToYAML:   *toYAML,
		ToSSH:    *toSSH,
		ToJSON:   *toJSON,
		Src:      *src,
		Dest:     *dest,
		ShowHelp: *showHelp,
	}
}

func CheckUseStdin(osStdinStat func() (fs.FileInfo, error)) bool {
	fi, err := osStdinStat()
	if err != nil {
		fmt.Println("Error getting stdin stat:", err)
		return false
	}
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		return true
	}
	return false
}

func CheckConvertArgvValid(args Args) (result bool, desc string) {
	trueCount := 0
	if args.ToJSON {
		trueCount++
	}
	if args.ToSSH {
		trueCount++
	}
	if args.ToYAML {
		trueCount++
	}

	if trueCount != 1 {
		return false, "Please specify either -to-yaml or -to-ssh or -to-json"
	}

	return true, ""
}

func CheckIOArgvValid(args Args) (result bool, desc string) {
	if args.Src == "" {
		return false, "Please specify source and destination file path"
	}

	// Check if src exists
	_, err := os.Stat(args.Src)
	if os.IsNotExist(err) {
		return false, fmt.Sprintf("Error: Source path '%s' does not exist", args.Src)
	}

	// Check if dist exists
	_, err = os.Stat(args.Dest)
	if os.IsNotExist(err) {
		// If dist doesn't exist, check if its parent directory exists
		parentDir := filepath.Dir(args.Dest)
		_, err := os.Stat(parentDir)
		if os.IsNotExist(err) {
			return false, fmt.Sprintf("Error: Parent directory of destination '%s' does not exist", args.Dest)
		}
	}

	return true, ""
}
