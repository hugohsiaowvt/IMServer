package utilservice

import (
	"net/http"
	"rz/restfulapi"
	"github.com/gin-gonic/gin"
)

type Country struct {
	Zone string `json:"zone"`
	Name string `json:"name"`
}

func Countrys() []*Country  {


	return []*Country{
		&Country{
			Zone: "0086",
			Name: "中国",
		},
		&Country{
			Zone: "001",
			Name: "美国",
		},
		&Country{
			Zone: "00853",
			Name: "澳门",
		},
		&Country{
			Zone: "001",
			Name: "加拿大",
		},
		&Country{
			Zone: "007",
			Name: "哈萨克斯坦",
		},
		&Country{
			Zone: "00998",
			Name: "乌兹别克斯坦",
		},
		&Country{
			Zone: "00996",
			Name: "吉尔吉斯斯坦",
		},
		&Country{
			Zone: "0090",
			Name: "土耳其",
		},
		&Country{
			Zone: "0033",
			Name: "法国",
		},
		&Country{
			Zone: "0049",
			Name: "德国",
		},
		&Country{
			Zone: "0044",
			Name: "英国",
		},
		&Country{
			Zone: "0039",
			Name: "意大利",
		},
		&Country{
			Zone: "00886",
			Name: "台湾",
		},
		&Country{
			Zone: "0060",
			Name: "马来西亚",
		},
		&Country{
			Zone: "0062",
			Name: "印度尼西亚",
		},
		&Country{
			Zone: "0061",
			Name: "澳大利亚",
		},
		&Country{
			Zone: "0064",
			Name: "新西兰",
		},
		&Country{
			Zone: "0063",
			Name: "菲律宾",
		},
		&Country{
			Zone: "0065",
			Name: "新加坡",
		},
		&Country{
			Zone: "0066",
			Name: "泰国",
		},
		&Country{
			Zone: "00673",
			Name: "文莱",
		},
		&Country{
			Zone: "0081",
			Name: "日本",
		},
		&Country{
			Zone: "0082",
			Name: "韩国",
		},
		&Country{
			Zone: "0084",
			Name: "越南",
		},
		&Country{
			Zone: "00852",
			Name: "香港",
		},
		&Country{
			Zone: "00855",
			Name: "柬埔寨",
		},
		&Country{
			Zone: "00856",
			Name: "老挝",
		},
		&Country{
			Zone: "00880",
			Name: "孟加拉国",
		},
		&Country{
			Zone: "0091",
			Name: "印度",
		},
		&Country{
			Zone: "0094",
			Name: "斯里兰卡",
		},
		&Country{
			Zone: "0095",
			Name: "缅甸",
		},
		&Country{
			Zone: "00960",
			Name: "马尔代夫",
		},
		&Country{
			Zone: "00976",
			Name: "蒙古",
		},
		&Country{
			Zone: "00975",
			Name: "不丹",
		},
		&Country{
			Zone: "007",
			Name: "俄罗斯",
		},
		&Country{
			Zone: "0030",
			Name: "希腊",
		},
		&Country{
			Zone: "0031",
			Name: "荷兰",
		},
		&Country{
			Zone: "0034",
			Name: "西班牙",
		},
		&Country{
			Zone: "00351",
			Name: "葡萄牙",
		},
		&Country{
			Zone: "00353",
			Name: "爱尔兰",
		},
		&Country{
			Zone: "0041",
			Name: "瑞士",
		},
		&Country{
			Zone: "0045",
			Name: "丹麦",
		},
		&Country{
			Zone: "0046",
			Name: "瑞典",
		},
		&Country{
			Zone: "0047",
			Name: "挪威",
		},
		&Country{
			Zone: "0055",
			Name: "巴西",
		},

	}
}

func GetCountrys(c *gin.Context) {
	c.JSON(http.StatusOK, restfulapi.Success(Countrys()))
}