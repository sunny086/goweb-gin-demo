package wt

import (
	"encoding/json"
	"goweb-gin-demo/global"
	"goweb-gin-demo/model/common/request"
	"goweb-gin-demo/model/wt"
	wtReq "goweb-gin-demo/model/wt/request"
	wtRes "goweb-gin-demo/model/wt/response"
)

type WtRuleService struct {
}

// CreateWtRule 创建WtRule记录
func (wtRuleService *WtRuleService) CreateWtRule(ruleRes wtReq.WtRuleRes) (err error) {
	rule := voToRule(ruleRes)
	err = global.GLOBAL_DB.Create(&rule).Error
	return err
}

// DeleteWtRuleByIds 批量删除WtRule记录
func (wtRuleService *WtRuleService) DeleteWtRuleByIds(ids request.IdsReq) (err error) {
	err = global.GLOBAL_DB.Delete(&[]wt.WtRule{}, "id in ?", ids.Ids).Error
	return err
}

// UpdateWtRule 更新WtRule记录
func (wtRuleService *WtRuleService) UpdateWtRule(ruleRes wtReq.WtRuleRes) (err error) {
	rule := voToRule(ruleRes)
	err = global.GLOBAL_DB.Updates(&rule).Error
	return err
}

// GetWtRule 根据id获取WtRule记录
func (wtRuleService *WtRuleService) GetWtRule(id uint) (err error, wtRule wt.WtRule) {
	err = global.GLOBAL_DB.Where("id = ?", id).First(&wtRule).Error
	return
}

// GetWtRuleInfoList 分页获取WtRule记录
func (wtRuleService *WtRuleService) GetWtRuleInfoList(info wtReq.WtRuleSearch) (err error, list []wtRes.WtRuleResult, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GLOBAL_DB.Model(&wt.WtRule{})
	var wtRules []wt.WtRule
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if info.Page <= 0 {
		limit = int(total)
	}

	err = db.Limit(limit).Offset(offset).Find(&wtRules).Error

	ruleResults := rulesToVOs(wtRules)

	return err, ruleResults, total
}

//数据转换一下, 需要把json数据转换为字符串
func voToRule(ruleRes wtReq.WtRuleRes) wt.WtRule {
	contentJson, _ := json.Marshal(ruleRes.Reporters)

	rule := wt.WtRule{
		GLOBAL_MODEL: global.GLOBAL_MODEL{ID: ruleRes.ID},
		UserId:       ruleRes.UserId,
		Reporters:    string(contentJson),
		StartWeek:    ruleRes.StartWeek,
		StartHour:    ruleRes.StartHour,
		EndWeek:      ruleRes.EndWeek,
		EndHour:      ruleRes.EndHour,
	}
	return rule
}

// 批量转换 数据转换, 把字符串转换为json
func rulesToVOs(rules []wt.WtRule) []wtRes.WtRuleResult {
	var ruleResults []wtRes.WtRuleResult
	for _, rule := range rules {
		ruleResult := ruleToResult(rule)
		ruleResults = append(ruleResults, ruleResult)
	}
	return ruleResults
}

//单个转换
func ruleToResult(rule wt.WtRule) wtRes.WtRuleResult {
	ruleResult := wtRes.WtRuleResult{}
	ruleResult.GLOBAL_MODEL = rule.GLOBAL_MODEL
	ruleResult.UserId = rule.UserId
	ruleResult.StartWeek = rule.StartWeek
	ruleResult.StartHour = rule.StartHour
	ruleResult.EndWeek = rule.EndWeek
	ruleResult.EndHour = rule.EndHour

	json.Unmarshal([]byte(rule.Reporters), &ruleResult.Reporters)
	return ruleResult
}
