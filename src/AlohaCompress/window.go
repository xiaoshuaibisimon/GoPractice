package main

import (
	"archive/zip"
	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
	"golang.org/x/exp/errors/fmt"
	"io"
	"os"
)

type Window interface {
	ShowWindow()
}

type ComWindow struct {
	Window
	*walk.MainWindow
}

func (c_w *ComWindow) OpenFileManager() string {
	dlg := new(walk.FileDialog)
	dlg.Title = "选择文件"
	dlg.Filter = "所有文件(*.*)|*.*|文本文档(*.txt)|*.txt"
	b, err := dlg.ShowOpen(c_w)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	fmt.Println(b)
	return dlg.FilePath
}

func (c_w *ComWindow) OpenDirManager() string {
	dlg := new(walk.FileDialog)
	dlg.Title = "选择路径"
	b, err := dlg.ShowBrowseFolder(c_w)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	fmt.Println(b)
	return dlg.FilePath
}

var (
	InfoLabel *walk.Label
	LabelText string
)

func (c_w *ComWindow) ShowWindow() {
	comWin := new(ComWindow)
	var unzipEdit *walk.LineEdit
	var saveUnzipEdit *walk.LineEdit
	var zipEdit *walk.LineEdit
	var saveZipEdit *walk.LineEdit

	var unzipBtn *walk.PushButton
	var saveUnzipBtn *walk.PushButton
	var zipBtn *walk.PushButton
	var saveZipBtn *walk.PushButton

	var startUnzipBtn *walk.PushButton
	var startZipBtn *walk.PushButton
	err := declarative.MainWindow{
		AssignTo: &comWin.MainWindow,
		Title:    "压缩工具",
		MinSize:  declarative.Size{480, 230},
		Layout:   declarative.HBox{},
		Children: []declarative.Widget{
			declarative.Composite{
				Layout: declarative.Grid{Columns: 2, Spacing: 10},
				Children: []declarative.Widget{
					declarative.LineEdit{
						Text:     "Input the ZIP file path",
						AssignTo: &unzipEdit,
					},
					declarative.PushButton{
						AssignTo: &unzipBtn,
						Text:     "Select The ZIP File",
						OnClicked: func() {
							filePath := comWin.OpenFileManager()
							unzipEdit.SetText(filePath)
						},
					},
					declarative.LineEdit{
						Text:     "Input the save path after unzip",
						AssignTo: &saveUnzipEdit,
					},
					declarative.PushButton{
						AssignTo: &saveUnzipBtn,
						Text:     "SavePath After Unzip",
						OnClicked: func() {
							filePath := comWin.OpenDirManager()
							saveUnzipEdit.SetText(filePath)
						},
					},

					declarative.LineEdit{
						Text:     "Input the file path which need to be compress",
						AssignTo: &zipEdit,
					},
					declarative.PushButton{
						AssignTo: &zipBtn,
						Text:     "Select The Common File",
						OnClicked: func() {
							filePath := comWin.OpenFileManager()
							zipEdit.SetText(filePath)
						},
					},
					declarative.LineEdit{
						Text:     "Input the save path after zip",
						AssignTo: &saveZipEdit,
					},
					declarative.PushButton{
						AssignTo: &saveZipBtn,
						Text:     "SavePath After Zip",
						OnClicked: func() {
							filePath := comWin.OpenDirManager()
							saveZipEdit.SetText(filePath)
						},
					},
					declarative.Label{
						AssignTo: &InfoLabel,
						Text:     "",
					},
				},
			},
			declarative.Composite{
				Layout: declarative.Grid{Rows: 2, Spacing: 40},
				Children: []declarative.Widget{
					declarative.PushButton{
						AssignTo: &startUnzipBtn,
						Text:     "Start Unzip",
						OnClicked: func() {
							if comWin.StartUnZip(unzipEdit.Text(), saveUnzipEdit.Text()) {
								LabelText = "Unzip successful"
							} else {
								LabelText = "Unzip failure"
							}
							Show("lab_window")
						},
					},
					declarative.PushButton{
						AssignTo: &startZipBtn,
						Text:     "Start Zip",
						OnClicked: func() {
							if comWin.StartZip(zipEdit.Text(), saveZipEdit.Text()) {
								LabelText = "Zip successful"
							} else {
								LabelText = "Zip failure"
							}
							Show("lab_window")
						},
					},
				},
			},
		},
	}.Create()

	if err != nil {
		fmt.Println(err)
	}

	comWin.SetX(600)
	comWin.SetY(800)
	comWin.Run()

}

type LabWindow struct {
	Window
}

func (l_w *LabWindow) ShowWindow() {
	InfoLabel.SetText(LabelText)
}

func Show(w_type string) {
	var Win Window
	switch w_type {
	case "main_window":
		Win = &ComWindow{}
	case "lab_window":
		Win = &LabWindow{}
	default:
		fmt.Println("Invalid window type parameter")
	}

	Win.ShowWindow()
}

func (c_w *ComWindow) StartUnZip(filePath string, savePath string) bool {
	rc, err := zip.OpenReader(filePath)
	defer rc.Close()
	if err != nil {
		fmt.Println(err)
		return false
	}

	for _, file := range rc.File {
		//newFilePath, err := UTF8ToGBK(savePath + "/" + file.Name)
		//if err != nil {
		//	fmt.Println(err)
		//	return false
		//}
		newFilePath := (savePath + "/" + file.Name)
		if file.FileInfo().IsDir() {
			err = os.MkdirAll(newFilePath, os.ModePerm)
			if err != nil {
				fmt.Println(err)
				return false
			}
		} else {
			srcFile, err := file.Open()
			defer srcFile.Close()
			if err != nil {
				fmt.Println(err)
				return false
			}

			dstFile, err := os.Create(newFilePath)
			defer dstFile.Close()
			if err != nil {
				fmt.Println(err)
				return false
			}

			_, err = io.Copy(dstFile, srcFile)
			if err != nil {
				fmt.Println(err)
				return false
			}
		}
	}
	return true
}

func (c_w *ComWindow) StartZip(filePath string, savePath string) bool {

	srcFile, err := os.Open(filePath)
	defer srcFile.Close()
	if err != nil {
		fmt.Println(err)
		return false
	}
	info, err := srcFile.Stat()
	if err != nil {
		fmt.Println(err)
		return false
	}
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		fmt.Println(err)
		return false
	}

	//newFilePath, err := UTF8ToGBK(savePath + "/" + info.Name() + ".zip")
	//if err != nil {
	//	fmt.Println(err)
	//	return false
	//}
	newFilePath := (savePath + "/" + info.Name() + ".zip")
	dstFile, err := os.Create(newFilePath)
	defer dstFile.Close()
	if err != nil {
		fmt.Println(err)
		return false
	}

	wr := zip.NewWriter(dstFile)
	defer wr.Close()
	writer, err := wr.CreateHeader(header)
	if err != nil {
		fmt.Println(err)
		return false
	}
	_, err = io.Copy(writer, srcFile)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
