package datacube

import (
	"encoding/json"
	"fmt"
	"github.com/silenceper/wechat/v2/util"
	"net/url"
	"strconv"
)

type AdSlot string

const (
	//广告位类型名称（ad_slot）	广告位类型
	SlotIdBizBottom         AdSlot = "SLOT_ID_BIZ_BOTTOM"         // 公众号底部广告
	SlotIdBizMidContext     AdSlot = "SLOT_ID_BIZ_MID_CONTEXT"    // 公众号文中广告
	SlotIdBizVideoEnd       AdSlot = "SLOT_ID_BIZ_VIDEO_END"      // 公众号视频后贴
	SlotIdBizSponsor        AdSlot = "SLOT_ID_BIZ_SPONSOR"        // 公众号互选广告
	SlotIdBizCps            AdSlot = "SLOT_ID_BIZ_CPS"            // 公众号返佣商品
	SlotIdWeappBanner       AdSlot = "SLOT_ID_WEAPP_BANNER"       // 小程序banner
	SlotIdWeappRewardVideo  AdSlot = "SLOT_ID_WEAPP_REWARD_VIDEO" // 小程序激励视频
	SlotIdWeappInterstitial AdSlot = "SLOT_ID_WEAPP_INTERSTITIAL" // 小程序插屏广告
	SlotIdWeappVideoFeeds   AdSlot = "SLOT_ID_WEAPP_VIDEO_FEEDS"  // 小程序视频广告
	SlotIdWeappVideoBegin   AdSlot = "SLOT_ID_WEAPP_VIDEO_BEGIN"  // 小程序视频前贴
	SlotIdWeappBox          AdSlot = "SLOT_ID_WEAPP_BOX"          // 小程序格子广告
)

const (
	publisherURL = "https://api.weixin.qq.com/publisher/stat"
)

const (
	actionPublisherAdPosGeneral = "publisher_adpos_general"
	actionPublisherCpsGeneral   = "publisher_cps_general"
	actionPublisherSettlement   = "publisher_settlement"
)

type BaseResp struct {
	ErrMsg string `json:"err_msg"`
	Ret    int    `json:"ret"`
}

//ResPublisherAdPos 公众号分广告位数据响应
type ResPublisherAdPos struct {
	util.CommonError
	BaseResp

	Base     BaseResp        `json:"base_resp"`
	List     []ResAdPosList  `json:"list"`
	Summary  ResAdPosSummary `json:"summary"`
	TotalNum int             `json:"total_num"`
}

type ResAdPosList struct {
	SlotID        int64   `json:"slot_id"`
	AdSlot        string  `json:"ad_slot"`
	Date          string  `json:"date"`
	ReqSuccCount  int     `json:"req_succ_count"`
	ExposureCount int     `json:"exposure_count"`
	ExposureRate  float64 `json:"exposure_rate"`
	ClickCount    int     `json:"click_count"`
	ClickRate     float64 `json:"click_rate"`
	Income        int     `json:"income"`
	Ecpm          float64 `json:"ecpm"`
}

type ResAdPosSummary struct {
	ReqSuccCount  int     `json:"req_succ_count"`
	ExposureCount int     `json:"exposure_count"`
	ExposureRate  float64 `json:"exposure_rate"`
	ClickCount    int     `json:"click_count"`
	ClickRate     float64 `json:"click_rate"`
	Income        int     `json:"income"`
	Ecpm          float64 `json:"ecpm"`
}

//ResPublisherCps 公众号返佣商品数据响应
type ResPublisherCps struct {
	util.CommonError
	BaseResp

	Base     BaseResp      `json:"base_resp"`
	List     []ResCpsList  `json:"list"`
	Summary  ResCpsSummary `json:"summary"`
	TotalNum int           `json:"total_num"`
}

type ResCpsList struct {
	Date            string  `json:"date"`
	ExposureCount   int     `json:"exposure_count"`
	ClickCount      int     `json:"click_count"`
	ClickRate       float64 `json:"click_rate"`
	OrderCount      int     `json:"order_count"`
	OrderRate       float64 `json:"order_rate"`
	TotalFee        int     `json:"total_fee"`
	TotalCommission int     `json:"total_commission"`
}

type ResCpsSummary struct {
	ExposureCount   int     `json:"exposure_count"`
	ClickCount      int     `json:"click_count"`
	ClickRate       float64 `json:"click_rate"`
	OrderCount      int     `json:"order_count"`
	OrderRate       float64 `json:"order_rate"`
	TotalFee        int     `json:"total_fee"`
	TotalCommission int     `json:"total_commission"`
}

//ResPublisherSettlement 公众号结算收入数据及结算主体信息响应
type ResPublisherSettlement struct {
	util.CommonError
	BaseResp

	Base              BaseResp         `json:"base_resp"`
	Body              string           `json:"body"`
	PenaltyAll        int              `json:"penalty_all"`
	RevenueAll        int64            `json:"revenue_all"`
	SettledRevenueAll int64            `json:"settled_revenue_all"`
	SettlementList    []SettlementList `json:"settlement_list"`
	TotalNum          int              `json:"total_num"`
}

type SettlementList struct {
	Date           string        `json:"date"`
	Zone           string        `json:"zone"`
	Month          string        `json:"month"`
	Order          int           `json:"order"`
	SettStatus     int           `json:"sett_status"`
	SettledRevenue int           `json:"settled_revenue"`
	SettNo         string        `json:"sett_no"`
	MailSendCnt    string        `json:"mail_send_cnt"`
	SlotRevenue    []SlotRevenue `json:"slot_revenue"`
}

type SlotRevenue struct {
	SlotID             string `json:"slot_id"`
	SlotSettledRevenue int    `json:"slot_settled_revenue"`
}

type ParamsPublisher struct {
	Action    string `json:"action"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Page      int    `json:"page"`
	PageSize  int    `json:"page_size"`
	AdSlot    AdSlot `json:"ad_slot"`
}

// fetchData 拉取统计数据
func (cube *DataCube) fetchData(params ParamsPublisher) (response []byte, err error) {
	accessToken, err := cube.GetAccessToken()
	if err != nil {
		return
	}

	v := url.Values{}
	v.Add("action", params.Action)
	v.Add("access_token", accessToken)
	v.Add("page", strconv.Itoa(params.Page))
	v.Add("page_size", strconv.Itoa(params.PageSize))
	v.Add("start_date", params.StartDate)
	v.Add("end_date", params.EndDate)
	if params.AdSlot != "" {
		v.Add("ad_slot", string(params.AdSlot))
	}

	uri := fmt.Sprintf("%s?%s", publisherURL, v.Encode())

	response, err = util.HTTPGet(uri)
	if err != nil {
		return
	}
	return
}

//GetPublisherAdPosGeneral 获取公众号分广告位数据
func (cube *DataCube) GetPublisherAdPosGeneral(startDate, endDate string, page, pageSize int, adSlot AdSlot) (resPublisherAdPos ResPublisherAdPos, err error) {
	params := ParamsPublisher{
		Action:    actionPublisherAdPosGeneral,
		StartDate: startDate,
		EndDate:   endDate,
		Page:      page,
		PageSize:  pageSize,
		AdSlot:    adSlot,
	}

	response, err := cube.fetchData(params)
	if err != nil {
		return
	}

	err = json.Unmarshal(response, &resPublisherAdPos)
	if err != nil {
		return
	}

	if resPublisherAdPos.CommonError.ErrCode != 0 {
		err = fmt.Errorf("GetPublisherAdPosGeneral Error , errcode=%d , errmsg=%s", resPublisherAdPos.CommonError.ErrCode, resPublisherAdPos.CommonError.ErrMsg)
		return
	}

	if resPublisherAdPos.BaseResp.Ret != 0 {
		err = fmt.Errorf("GetPublisherAdPosGeneral Error , errcode=%d , errmsg=%s", resPublisherAdPos.BaseResp.Ret, resPublisherAdPos.BaseResp.ErrMsg)
		return
	}

	if resPublisherAdPos.Base.Ret != 0 {
		err = fmt.Errorf("GetPublisherAdPosGeneral Error , errcode=%d , errmsg=%s", resPublisherAdPos.Base.Ret, resPublisherAdPos.Base.ErrMsg)
		return
	}
	return
}

//GetPublisherCpsGeneral 获取公众号返佣商品数据
func (cube *DataCube) GetPublisherCpsGeneral(startDate, endDate string, page, pageSize int) (resPublisherCps ResPublisherCps, err error) {
	params := ParamsPublisher{
		Action:    actionPublisherCpsGeneral,
		StartDate: startDate,
		EndDate:   endDate,
		Page:      page,
		PageSize:  pageSize,
	}

	response, err := cube.fetchData(params)
	if err != nil {
		return
	}

	err = json.Unmarshal(response, &resPublisherCps)
	if err != nil {
		return
	}

	if resPublisherCps.CommonError.ErrCode != 0 {
		err = fmt.Errorf("GetPublisherCpsGeneral Error , errcode=%d , errmsg=%s", resPublisherCps.CommonError.ErrCode, resPublisherCps.CommonError.ErrMsg)
		return
	}

	if resPublisherCps.BaseResp.Ret != 0 {
		err = fmt.Errorf("GetPublisherCpsGeneral Error , errcode=%d , errmsg=%s", resPublisherCps.BaseResp.Ret, resPublisherCps.BaseResp.ErrMsg)
		return
	}

	if resPublisherCps.Base.Ret != 0 {
		err = fmt.Errorf("GetPublisherCpsGeneral Error , errcode=%d , errmsg=%s", resPublisherCps.Base.Ret, resPublisherCps.Base.ErrMsg)
		return
	}

	return
}

//GetPublisherSettlement 获取公众号结算收入数据及结算主体信息
func (cube *DataCube) GetPublisherSettlement(startDate, endDate string, page, pageSize int) (resPublisherSettlement ResPublisherSettlement, err error) {
	params := ParamsPublisher{
		Action:    actionPublisherSettlement,
		StartDate: startDate,
		EndDate:   endDate,
		Page:      page,
		PageSize:  pageSize,
	}

	response, err := cube.fetchData(params)
	if err != nil {
		return
	}

	err = json.Unmarshal(response, &resPublisherSettlement)
	if err != nil {
		return
	}

	if resPublisherSettlement.CommonError.ErrCode != 0 {
		err = fmt.Errorf("GetPublisherSettlement Error , errcode=%d , errmsg=%s", resPublisherSettlement.CommonError.ErrCode, resPublisherSettlement.CommonError.ErrMsg)
		return
	}

	if resPublisherSettlement.BaseResp.Ret != 0 {
		err = fmt.Errorf("GetPublisherSettlement Error , errcode=%d , errmsg=%s", resPublisherSettlement.BaseResp.Ret, resPublisherSettlement.BaseResp.ErrMsg)
		return
	}

	if resPublisherSettlement.Base.Ret != 0 {
		err = fmt.Errorf("GetPublisherSettlement Error , errcode=%d , errmsg=%s", resPublisherSettlement.Base.Ret, resPublisherSettlement.Base.ErrMsg)
		return
	}
	return
}
