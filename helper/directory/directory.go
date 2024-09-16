package directory

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	cnf "github.com/akash-aman/tx/config"
)

/**
 * Get directory name.
 */
func GetDir(path string) string {

	cwd := GetCwd(path)

	name := filepath.Base(cwd)

	if name == "" {
		fmt.Printf("Error getting directory name: %v\n", name)
		os.Exit(1)
	}

	return name
}

/**
 * Get current working directory.
 */
func GetCwd(path string) string {

	if strings.HasPrefix(path, "~") {
		usr, err := user.Current()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		path = filepath.Join(usr.HomeDir, path[1:])
	}

	_, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	absPath, err := filepath.Abs(path)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return absPath
}

/**
 * Count files in a directory.
 */
func CountFiles(path string) int {
	count := 0

	cnf.AddConfig.Path = path
	cnf.AddConfig.Load()

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		for _, ignore := range cnf.AddConfig.Ignore {
			if strings.Contains(path, ignore) {
				return nil
			}
		}

		if info.IsDir() {
			return nil
		}

		count++

		return nil
	})

	if err != nil {
		fmt.Println("Error counting files:", err)
		os.Exit(1)
	}

	return count
}

/**
 * Copy file from source to destination.
 */
func CopyFile(src, dst string) error {

	err := os.MkdirAll(filepath.Dir(dst), 0755)
	if err != nil {
		return err
	}

	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}

	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}
