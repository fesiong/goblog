package manageController

import (
	"github.com/kataras/iris/v12"
	"kandaoni.com/anqicms/config"
	"kandaoni.com/anqicms/library"
	"kandaoni.com/anqicms/model"
	"kandaoni.com/anqicms/provider"
	"kandaoni.com/anqicms/request"
)

func PluginCommentList(ctx iris.Context) {
	currentSite := provider.CurrentSite(ctx)
	currentPage := ctx.URLParamIntDefault("current", 1)
	pageSize := ctx.URLParamIntDefault("pageSize", 20)
	comments, total, err := currentSite.GetCommentList(0, 0, "id desc", currentPage, pageSize, 0)
	if err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  "",
		})
		return
	}
	type miniArticle struct {
		Id    uint
		Title string
	}
	for i, v := range comments {
		var article miniArticle
		err = currentSite.DB.Model(&model.Archive{}).Where("id = ?", v.ArchiveId).Scan(&article).Error
		if err == nil {
			comments[i].ItemTitle = article.Title
		}
	}

	ctx.JSON(iris.Map{
		"code":  config.StatusOK,
		"msg":   "",
		"total": total,
		"data":  comments,
	})
}

func PluginCommentDetail(ctx iris.Context) {
	currentSite := provider.CurrentSite(ctx)
	id := uint(ctx.URLParamIntDefault("id", 0))
	comment, err := currentSite.GetCommentById(id)
	if err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	ctx.JSON(iris.Map{
		"code": config.StatusOK,
		"msg":  "",
		"data": comment,
	})
}

func PluginCommentDetailForm(ctx iris.Context) {
	currentSite := provider.CurrentSite(ctx)
	var req request.PluginComment
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	comment, err := currentSite.GetCommentById(req.Id)
	if err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	// 将单个&nbsp;替换为空格
	req.Content = library.ReplaceSingleSpace(req.Content)
	comment.UserName = req.UserName
	comment.Content = req.Content
	comment.Status = 1
	if req.Ip == "" {
		req.Ip = ctx.RemoteAddr()
	}
	comment.Ip = req.Ip

	err = comment.Save(currentSite.DB)
	if err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	currentSite.AddAdminLog(ctx, ctx.Tr("修改评论内容：%d", comment.Id))

	ctx.JSON(iris.Map{
		"code": config.StatusOK,
		"msg":  ctx.Tr("评论已更新"),
	})
}

func PluginCommentDelete(ctx iris.Context) {
	currentSite := provider.CurrentSite(ctx)
	var req request.PluginComment
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}
	comment, err := currentSite.GetCommentById(req.Id)
	if err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	err = comment.Delete(currentSite.DB)
	if err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	currentSite.AddAdminLog(ctx, ctx.Tr("修改文档模型：%d => %s", comment.Id, comment.Content))

	ctx.JSON(iris.Map{
		"code": config.StatusOK,
		"msg":  ctx.Tr("评论已删除"),
	})
}

// 处理审核状态
func PluginCommentCheck(ctx iris.Context) {
	currentSite := provider.CurrentSite(ctx)
	var req request.PluginComment
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	if len(req.Ids) > 0 {
		err := currentSite.DB.Model(&model.Comment{}).Where("`id` IN (?)", req.Ids).UpdateColumn("status", req.Status).Error
		if err != nil {
			ctx.JSON(iris.Map{
				"code": config.StatusFailed,
				"msg":  err.Error(),
			})
			return
		}

		currentSite.AddAdminLog(ctx, ctx.Tr("批量审核评论：%v", req.Ids))

		ctx.JSON(iris.Map{
			"code": config.StatusOK,
			"msg":  ctx.Tr("评论已更新"),
		})
		return
	}

	comment, err := currentSite.GetCommentById(req.Id)
	if err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	if comment.Status != model.StatusOk {
		comment.Status = model.StatusOk
	} else {
		comment.Status = model.StatusWait
	}
	err = comment.Save(currentSite.DB)
	if err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	currentSite.AddAdminLog(ctx, ctx.Tr("审核评论：%d", comment.Id))

	ctx.JSON(iris.Map{
		"code": config.StatusOK,
		"msg":  ctx.Tr("评论已更新"),
	})
}
