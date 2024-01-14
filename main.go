package main

import (
	_ "embed"
	"flag"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

//go:embed version
var version string

var videxts = []string{"mkv", "mov", "avi", "wmv", "m4v", "webm"}
var copyexts = []string{"mkv", "mov"}

var path string
var verbose bool
var recursive bool
var copy bool
var delete bool
var move bool
var moveto string
var ex string
var args []string
var exitonerror bool

var converted, failed int
var totaloldsize, totalnewsize int64

func main() {
	e, err := os.Executable()
	if err != nil {
		panic(err)
	}

	flag := flag.NewFlagSet("ConvertToMP4", flag.ExitOnError)
	flag.StringVar(&path, "p", filepath.Dir(e), "Path (defaults to current working directory)")
	flag.BoolVar(&verbose, "v", false, "Verbose")
	flag.BoolVar(&recursive, "r", false, "Recursive")
	flag.BoolVar(&copy, "c", false, "Copy the existing frames for formats that most likely support it")
	flag.BoolVar(&delete, "d", false, "Delete old files")
	flag.BoolVar(&move, "m", false, "Move old files")
	flag.StringVar(&moveto, "mt", "", "Folder to move old files to")
	flag.StringVar(&ex, "exec", "ffmpeg", "Overwrite FFmpeg executable")
	arg := *flag.String("args", "", "Args to pass to FFmpeg")
	flag.BoolVar(&exitonerror, "e", false, "Exit on FFmpeg error")
	ver := flag.Bool("version", false, "Get version info")
	flag.Parse(os.Args[1:])

	if *ver {
		fmt.Println("ConvertToMP4 version "+version, runtime.GOOS+"/"+runtime.GOARCH)
		os.Exit(0)
	}

	if path == "" {
		exit(1, "Path empty")
	}
	if !filepath.IsAbs(path) {
		p, err := filepath.Abs(path)
		if err != nil {
			exit(1, err.Error())
		}
		path = p
	}

	if delete && move {
		exit(1, "Can't delete and move old files")
	}

	if move {
		if moveto == "" {
			exit(1, "No path to move old files to specified")
		}
		if runtime.GOOS == "windows" {
			if path[0] != moveto[0] {
				exit(1, "Cross partition file moving is currently not supported on Windows")
			}
		}
	}

	if arg != "" {
		args = strings.Split(arg, " ")
	}

	process(path)

	if verbose {
		fmt.Println()
		fmt.Println()
		fmt.Println("Successfully converted " + strconv.Itoa(converted) + " files")
		fmt.Println("Failed to convert " + strconv.Itoa(failed) + " files")
		fmt.Println("Total size of old files: " + humanfilesize(totaloldsize))
		fmt.Println("Total size of converted files: " + humanfilesize(totalnewsize))
	}
}

func process(dir string) {
	e, err := os.ReadDir(dir)
	if err != nil {
		exit(1, err.Error())
	}
	for _, v := range e {
		n := v.Name()
		if v.IsDir() {
			if n == "." || n == ".." {
				continue
			}
			if recursive {
				process(filepath.Join(dir, n))
			}
			continue
		}
		if !strings.Contains(n, ".") {
			continue
		}
		split := strings.Split(n, ".")
		splitlen := len(split)
		var file, ext string
		if splitlen > 2 {
			file = strings.Join(split[:splitlen-1], ".")
			ext = split[splitlen-1]
		} else {
			file = split[0]
			ext = split[1]
		}
		if !instringarray(ext, videxts) {
			continue
		}

		fp := filepath.Join(dir, n)
		ofp := filepath.Join(dir, file+".mp4")
		if exists(ofp) {
			ofp = filepath.Join(dir, file+"_converttomp4.mp4")
			if exists(ofp) {
				fmt.Println("There already exists an converted file for " + fp)
				continue
			}
		}

		var a []string
		a = append(a, "-i", fp)
		a = append(a, args...)
		if copy {
			if instringarray(ext, copyexts) {
				a = append(a, "-c", "copy")
			}
		}
		a = append(a, ofp)
		if verbose {
			fmt.Println("Converting", fp)
		}
		_, err := exec.Command(ex, a...).Output()
		if err != nil {
			fmt.Println("Failed converting " + fp)
			outerr, ok := err.(*exec.ExitError)
			if ok {
				fmt.Println(string(outerr.Stderr))
			} else {
				fmt.Println(err)
			}
			os.Remove(ofp)
			if exitonerror {
				os.Exit(1)
			}
			failed++
			continue
		}
		converted++
		var oldsize, newsize int64
		old, err := os.Stat(fp)
		if err == nil {
			oldsize = old.Size()
			totaloldsize += oldsize
		}
		new, err := os.Stat(ofp)
		if err == nil {
			newsize = new.Size()
			totalnewsize += newsize
		}
		if delete {
			os.Remove(fp)
		}
		if move {
			mvdir := moveto + strings.TrimPrefix(dir, path)
			os.MkdirAll(mvdir, 0755)
			os.Rename(fp, filepath.Join(mvdir, n))
		}
		if verbose {
			fmt.Println("Converted", fp, "(old filesize: "+humanfilesize(oldsize)+" converted filesize: "+humanfilesize(newsize)+")")
		}
	}
}

func exit(code int, msg string) {
	if msg != "" {
		fmt.Println(msg)
	}
	os.Exit(code)
}

func instringarray(val string, array []string) bool {
	for i := range array {
		if array[i] == val {
			return true
		}
	}
	return false
}

func exists(fp string) bool {
	_, err := os.Stat(fp)
	return err == nil
}

var suffixes = []string{"B", "KB", "MB", "GB", "TB"}

func humanfilesize(size int64) string {
	if size == 0 {
		return "0B"
	}

	base := math.Log(float64(size)) / math.Log(1024)
	getSize := round(math.Pow(1024, base-math.Floor(base)), .5, 2)
	getSuffix := suffixes[int(math.Floor(base))]
	return strconv.FormatFloat(getSize, 'f', -1, 64) + " " + string(getSuffix)
}

func round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}
