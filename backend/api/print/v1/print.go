// =================================================================================
// 打印设置：全局单据打印样式（抬头 / 副标题 / 页脚 / 日期 / 电子章）
// =================================================================================

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// PrintSettingModel 打印设置模型
type PrintSettingModel struct {
	Title        string `json:"title"        dc:"抬头标题，空则打印时仅显示单据名"`
	SubTitle     string `json:"subTitle"     dc:"副标题（如公司英文名 / 部门名）"`
	Footer       string `json:"footer"       dc:"页脚文字"`
	ShowDate     bool   `json:"showDate"     dc:"是否在右上角显示打印日期"`
	ShowSeal     bool   `json:"showSeal"     dc:"是否加盖电子章"`
	SealImage    string `json:"sealImage"    dc:"印章图片相对路径"`
	SealImageUrl string `json:"sealImageUrl" dc:"印章图片绝对地址（仅查询返回）"`
}

// SettingGetReq 获取打印设置请求
type SettingGetReq struct {
	g.Meta `path:"/print/setting/get" method:"get" tags:"打印设置" summary:"获取打印设置"`
}

// SettingGetRes 获取打印设置响应
type SettingGetRes struct {
	Setting *PrintSettingModel `json:"setting" dc:"打印设置，未配置时返回默认值"`
}

// SettingSaveReq 保存打印设置请求
type SettingSaveReq struct {
	g.Meta `path:"/print/setting/save" method:"post" tags:"打印设置" summary:"保存打印设置"`

	Title     string `json:"title"     v:"max-length:120#抬头标题最多 120 个字符" dc:"抬头标题"`
	SubTitle  string `json:"subTitle"  v:"max-length:120#副标题最多 120 个字符" dc:"副标题"`
	Footer    string `json:"footer"    v:"max-length:255#页脚文字最多 255 个字符" dc:"页脚文字"`
	ShowDate  bool   `json:"showDate"  dc:"是否显示打印日期"`
	ShowSeal  bool   `json:"showSeal"  dc:"是否加盖电子章"`
	SealImage string `json:"sealImage" v:"max-length:255#印章图片地址过长" dc:"印章图片相对路径"`
}

// SettingSaveRes 保存打印设置响应
type SettingSaveRes struct{}
