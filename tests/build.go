/**
 * @Time : 2024/5/20 14:06
 * @File : build
 * @Software: GoLand
 * @Description: 支持打包操作，打包成功将配置文件和静态文件拷贝到对应目录下
 * @Version: 1.0.0
 */
package main

import (
	"archive/zip"
	"fmt"
	"go-wx-download/pkg/utils"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// copyDir 复制源目录到目标目录
func copyDir(src, dst string) error {
	// 打开源目录
	srcDir, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcDir.Close()
	// 遍历源目录中的文件和子目录
	entries, err := srcDir.Readdir(-1)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		// 如果是目录，递归调用copyDir
		if entry.IsDir() {
			if err := os.MkdirAll(dstPath, os.ModePerm); err != nil {
				return err
			}
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// 如果是文件，复制文件
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}

// copyFile 将源文件复制到目标文件
func copyFile(src, dst string) error {
	// 打开源文件和目标文件
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// 复制文件内容
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}
	return nil
}

func copyResource() {
	// 获取当前工作路径
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("获取当前工作路径失败: %s\n", err)
		return
	}
	fmt.Println("cwd：", cwd)
	dir, err := os.ReadDir("./")
	if err != nil {
		log.Fatalln("读取项目根目录失败", err)
	}
	// 存放根目录下的文件夹
	var folder []string
	for _, entry := range dir {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		if info.IsDir() && strings.HasPrefix(entry.Name(), "wxdown-") {
			folder = append(folder, entry.Name())
		}
	}
	// 遍历输出文件夹
	for _, f := range folder {
		fmt.Println(f)
		web := filepath.Join(cwd, f, "web")
		utils.IsNotExistCreate(web)
		if err := copyDir("web", web); err != nil {
			fmt.Printf("Error copying directory: %v\n", err)
		} else {
			fmt.Println("Copy complete")
		}
		// 拼接路径
		dst := filepath.Join(f, "config.yaml")
		utils.CopyFile(dst, "./config.yaml")
	}
	fmt.Println("完成资源拷贝!....")
	fmt.Println("打包!....")
	// 遍历输出文件夹
	for _, f := range folder {
		// 5、打包
		sprintf := fmt.Sprintf("%s.zip", f)
		fmt.Println(sprintf)

		outputFile, err := os.Create(sprintf)
		if err != nil {
			log.Fatal(err)
		}
		defer outputFile.Close()

		zipWriter := zip.NewWriter(outputFile)
		defer zipWriter.Close()

		// Replace "directory_to_zip" with the actual directory you want to zip
		directoryToZip := f
		baseDir := filepath.Dir(directoryToZip)

		err = zipDirectory(zipWriter, directoryToZip, baseDir)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// addFileToZip adds a file to the provided zip.Writer
func addFileToZip(zipWriter *zip.Writer, filename string, baseDir string) error {
	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	// Create a relative path for the file in the ZIP archive
	relativePath, err := filepath.Rel(baseDir, filename)
	if err != nil {
		return err
	}
	relativePath = filepath.ToSlash(relativePath) // Ensure the path uses '/' as separator

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = relativePath

	if info.IsDir() {
		header.Name += "/"
	} else {
		header.Method = zip.Deflate
	}

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}

	_, err = io.Copy(writer, fileToZip)
	return err
}

// zipDirectory recursively adds a directory to the provided zip.Writer
func zipDirectory(zipWriter *zip.Writer, dir string, baseDir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return addFileToZip(zipWriter, path, baseDir)
	})
}
func main() {
	copyResource()
}
