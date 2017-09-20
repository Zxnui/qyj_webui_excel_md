package main

import (
	"github.com/andlabs/ui"
	"crypto/md5"
	"github.com/tealeg/xlsx"
	"fmt"
	"time"
	"strconv"
)

func main() {
	err:=ui.Main(func() {

		/************************公共资源***********************/
		fileFrom := "nil"//上传文件
		fileTemp := ""//上传文件
		windowMsg := ui.NewWindow("提示",100,100,false)
		windowFile := ui.NewWindow("下载",1,1,false)
		status:=0

		/************************公共资源***********************/

		/************************样式***********************/
		tab :=ui.NewTab()
		box1 := ui.NewVerticalBox()
		g1 := ui.NewGroup("his模板")
		button1 := ui.NewButton("点击下载模本")
		g1.SetChild(button1)
		box1.Append(g1,false)

		g2 := ui.NewGroup("pacs模板")
		buttonP := ui.NewButton("点击下载模本")
		g2.SetChild(buttonP)
		box1.Append(g2,false)
		tab.InsertAt("模板下载",0,box1)

		box2 := ui.NewVerticalBox()

		group1 := ui.NewGroup("第一步")
		box2_1 := ui.NewVerticalBox()
		up2 := ui.NewButton("上传模本")
		box2_1.Append(up2,false)

		box2_2 := ui.NewHorizontalBox()
		label2_1 := ui.NewLabel("上传进度条:")
		vLab2_1 := ui.NewLabel("")
		box2_2.Append(label2_1,false)
		box2_2.Append(vLab2_1,false)
		box2_1.Append(box2_2,false)
		group1.SetChild(box2_1)

		box2.Append(ui.NewLabel(""),false)
		box2.Append(group1,false)

		box2_3 := ui.NewHorizontalBox()
		box2_4 := ui.NewVerticalBox()
		group2 := ui.NewGroup("第二步")
		button2 := ui.NewButton("加密数据")
		box2_4.Append(button2,false)

		label2_2 := ui.NewLabel("加密进度条:")
		vLab2_2 :=ui.NewLabel("")
		box2_3.Append(label2_2,false)
		box2_3.Append(vLab2_2,false)
		box2_4.Append(box2_3,false)
		group2.SetChild(box2_4)
		box2.Append(ui.NewLabel(""),false)
		box2.Append(group2,false)


		group3 := ui.NewGroup("第三步")
		box2_5 := ui.NewVerticalBox()

		label3 := ui.NewLabel("加密完成后，才能下载加密后的数据:")
		box2_5.Append(label3,false)
		button3 := ui.NewButton("点击下载加密后数据")
		box2_5.Append(button3,false)

		box2_6:=ui.NewHorizontalBox()
		label2_6 := ui.NewLabel("下载进度条:")
		vLab2_6 :=ui.NewLabel("")
		box2_6.Append(label2_6,false)
		box2_6.Append(vLab2_6,false)
		box2_5.Append(box2_6,false)

		group3.SetChild(box2_5)
		box2.Append(ui.NewLabel(""),false)
		box2.Append(group3,false)

		tab.InsertAt("数据加密",1,box2)

		window := ui.NewWindow("数据加密", 600, 400, false)
		window.SetChild(tab)

		/************************样式***********************/

		/************************动态处理***********************/
		//点击上传数据
		up2.OnClicked(func(*ui.Button) {
			fileFrom = ui.OpenFile(ui.NewWindow("打开",400,300,false))
			if len(fileFrom)<5 {
				fileFrom = "nil"
				return
			}
			vLab2_1.SetText( "文件上传中，请稍后......")
			fileTemp = "./temp/temp"+strconv.FormatInt(time.Now().Unix(),10)+".xlsx"
			_,err:=copyFile(fileTemp,fileFrom)
			check(err)
			vLab2_1.SetText( "文件上传成功")
			ui.MsgBox(windowMsg,"提示","数据上传成功")
		})

		//下载加密后的数据
		button3.OnClicked(func(*ui.Button) {
			if fileFrom!="nil" {
				if status==0{
					ui.MsgBoxError(windowMsg,"提示","请先对数据进行加密")
				}else if status==1 {
					ui.MsgBoxError(windowMsg,"提示","数据加密中，请稍后......")
				}else{
					fileAddress:=ui.SaveFile(windowFile)
					vLab2_6.SetText( "文件下载中，请稍后......")
					copyFile(fileAddress+".xlsx",fileTemp)
					vLab2_6.SetText( "文件下载成功")
					ui.MsgBox(windowMsg,"提示","加密后的数据下载成功")
				}
			}else{
				ui.MsgBoxError(windowMsg,"提示","请先上传需要加密的文件")
			}
		})

		//开始数据加密
		button2.OnClicked(func(*ui.Button) {
			if fileFrom!="nil" {
				vLab2_2.SetText("数据加密中，请稍后......")
				status = 1

				//读取文件，转码
				xlFile, err := xlsx.OpenFile(fileTemp)
				check(err)
				sheet := xlFile.Sheets[0]
				size:=0
				if len(sheet.Rows[0].Cells) == 9{
					size=2
				}else if len(sheet.Rows[0].Cells) == 20{
					size=6
				}else{
					ui.MsgBoxError(windowMsg,"提示","请使用我们给与的模板导入数据")
					return
				}

				//不同模板，数据加密单位不同
				for i, row := range sheet.Rows {
					if i>0 {
						b:=md5.Sum([]byte(row.Cells[size].Value))
						sheet.Rows[i].Cells[size].Value =fmt.Sprintf("%x", b)
					}
				}

				//转码完成，保存
				err=xlFile.Save(fileTemp)
				check(err)
				//状态重置
				status = 2
				vLab2_2.SetText("数据加密完成，请下载加密后数据")
				ui.MsgBox(windowMsg,"提示","数据加密完成，请下载加密后数据")
			}else{
				ui.MsgBoxError(windowMsg,"提示","请先上传需要加密的文件")
			}
		})

		//点击下载his模板
		button1.OnClicked(func(*ui.Button) {
			fileAddress:=ui.SaveFile(windowFile)
			if len(fileAddress)<5 {
				return
			}
			var filename = "./file/his.xlsx"
			_,err:=copyFile(fileAddress+".xlsx",filename)
			check(err)
			ui.MsgBox(windowMsg,"提示","His模板下载成功")
		})

		//点击下载pacs模板
		buttonP.OnClicked(func(*ui.Button) {
			fileAddress:=ui.SaveFile(windowFile)
			if len(fileAddress)<5 {
				return
			}
			var filename = "./file/pacs.xlsx"
			_,err:=copyFile(fileAddress+".xlsx",filename)
			check(err)
			ui.MsgBox(windowMsg,"提示","Pacs模板下载成功")
		})
		/************************动态处理***********************/

		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})
		window.Show()
	})

	if err != nil {
		panic(err)
	}
}