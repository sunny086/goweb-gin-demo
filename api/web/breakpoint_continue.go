package web

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"goweb-gin-demo/global"
	"goweb-gin-demo/model/common/response"
	webRes "goweb-gin-demo/model/web/response"
	"goweb-gin-demo/utils"
)

// @Tags FileUploadAndDownload
// @Summary 断点续传到服务器
// @Security ApiKeyAuth
// @accept multipart/form-data
// @Produce  application/json
// @Param file formData file true "an example for breakpoint resume, 断点续传示例"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"切片创建成功"}"
// @Router /fileUploadAndDownload/breakpointContinue [post]
func (u *FileUploadAndDownloadApi) BreakpointContinue(c *gin.Context) {
	fileMd5 := c.Request.FormValue("fileMd5")
	fileName := c.Request.FormValue("fileName")
	chunkMd5 := c.Request.FormValue("chunkMd5")
	chunkNumber, _ := strconv.Atoi(c.Request.FormValue("chunkNumber"))
	chunkTotal, _ := strconv.Atoi(c.Request.FormValue("chunkTotal"))
	_, FileHeader, err := c.Request.FormFile("file")
	if err != nil {
		global.GLOBAL_LOG.Error("接收文件失败!", zap.Any("err", err))
		response.FailWithMessage("接收文件失败", c)
		return
	}
	f, err := FileHeader.Open()
	if err != nil {
		global.GLOBAL_LOG.Error("文件读取失败!", zap.Any("err", err))
		response.FailWithMessage("文件读取失败", c)
		return
	}
	defer func(f multipart.File) {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(f)
	cen, _ := ioutil.ReadAll(f)
	if !utils.CheckMd5(cen, chunkMd5) {
		global.GLOBAL_LOG.Error("检查md5失败!", zap.Any("err", err))
		response.FailWithMessage("检查md5失败", c)
		return
	}
	err, file := fileUploadAndDownloadService.FindOrCreateFile(fileMd5, fileName, chunkTotal)
	if err != nil {
		global.GLOBAL_LOG.Error("查找或创建记录失败!", zap.Any("err", err))
		response.FailWithMessage("查找或创建记录失败", c)
		return
	}
	err, pathc := utils.BreakPointContinue(cen, fileName, chunkNumber, chunkTotal, fileMd5)
	if err != nil {
		global.GLOBAL_LOG.Error("断点续传失败!", zap.Any("err", err))
		response.FailWithMessage("断点续传失败", c)
		return
	}

	if err = fileUploadAndDownloadService.CreateFileChunk(file.ID, pathc, chunkNumber); err != nil {
		global.GLOBAL_LOG.Error("创建文件记录失败!", zap.Any("err", err))
		response.FailWithMessage("创建文件记录失败", c)
		return
	}
	response.OkWithMessage("切片创建成功", c)
}

// @Tags FileUploadAndDownload
// @Summary 查找文件
// @Security ApiKeyAuth
// @accept multipart/form-data
// @Produce  application/json
// @Param file formData file true "Find the file, 查找文件"qqq
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查找成功"}"
// @Router /fileUploadAndDownload/findFile [post]
func (u *FileUploadAndDownloadApi) FindFile(c *gin.Context) {
	fileMd5 := c.Query("fileMd5")
	fileName := c.Query("fileName")
	chunkTotal, _ := strconv.Atoi(c.Query("chunkTotal"))
	err, file := fileUploadAndDownloadService.FindOrCreateFile(fileMd5, fileName, chunkTotal)
	if err != nil {
		global.GLOBAL_LOG.Error("查找失败!", zap.Any("err", err))
		response.FailWithMessage("查找失败", c)
	} else {
		response.OkWithDetailed(webRes.FileResponse{File: file}, "查找成功", c)
	}
}

// @Tags FileUploadAndDownload
// @Summary 创建文件
// @Security ApiKeyAuth
// @accept multipart/form-data
// @Produce  application/json
// @Param file formData file true "上传文件完成"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"file uploaded, 文件创建成功"}"
// @Router /fileUploadAndDownload/breakpointContinueFinish [post]
func (b *FileUploadAndDownloadApi) BreakpointContinueFinish(c *gin.Context) {
	fileMd5 := c.Query("fileMd5")
	fileName := c.Query("fileName")
	err, filePath := utils.MakeFile(fileName, fileMd5)
	if err != nil {
		global.GLOBAL_LOG.Error("文件创建失败!", zap.Any("err", err))
		response.FailWithDetailed(webRes.FilePathResponse{FilePath: filePath}, "文件创建失败", c)
	} else {
		response.OkWithDetailed(webRes.FilePathResponse{FilePath: filePath}, "文件创建成功", c)
	}
}

// @Tags FileUploadAndDownload
// @Summary 删除切片
// @Security ApiKeyAuth
// @accept multipart/form-data
// @Produce  application/json
// @Param file formData file true "删除缓存切片"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"缓存切片删除成功"}"
// @Router /fileUploadAndDownload/removeChunk [post]
func (u *FileUploadAndDownloadApi) RemoveChunk(c *gin.Context) {
	fileMd5 := c.Query("fileMd5")
	fileName := c.Query("fileName")
	filePath := c.Query("filePath")
	err := utils.RemoveChunk(fileMd5)
	if err != nil {
		return
	}
	err = fileUploadAndDownloadService.DeleteFileChunk(fileMd5, fileName, filePath)
	if err != nil {
		global.GLOBAL_LOG.Error("缓存切片删除失败!", zap.Any("err", err))
		response.FailWithDetailed(webRes.FilePathResponse{FilePath: filePath}, "缓存切片删除失败", c)
	} else {
		response.OkWithDetailed(webRes.FilePathResponse{FilePath: filePath}, "缓存切片删除成功", c)
	}
}