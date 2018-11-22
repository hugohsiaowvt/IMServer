package fileservice

import (
	"encoding/json"
	"net/http"
	"rz/config"
	"rz/network"
	"rz/restfulapi"
	"github.com/gin-gonic/gin"
)

func DirAssign(c *gin.Context) {

	res := &AssignFidResponse{}

	resp, err := network.Get(config.FILE_SERVER_ADDRESS + "/dir/assign", nil, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, restfulapi.Error(1, err.Error()))
		return
	}
	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadRequest, restfulapi.Error(1, "獲取fid返回狀態有誤"))
		return
	}

	if err := json.Unmarshal([]byte(resp.Body), &res); err != nil {
		panic(err)
	}
	
	c.JSON(http.StatusOK, restfulapi.Success(res))
}
