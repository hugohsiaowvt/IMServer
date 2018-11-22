package restfulapi

import(
	"net/http"

	"github.com/gin-gonic/gin"
)

type RestfulAPIResponse struct {
	Status	int
	Data	interface{}
}

type SuccessMessage struct {
	Message		string
}

type ErrorMessage struct {
	ErrorCode	int8
	ErrorMsg	string
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, RestfulAPIResponse{0, data})
}

func SuccessMsg(c *gin.Context, msg string) {
	res := SuccessMessage {
		Message:	msg,
	}
	c.JSON(http.StatusOK, RestfulAPIResponse{0, res})
}

func Error(c *gin.Context, code int8, msg string) {
	err := ErrorMessage {
		ErrorCode:	code,
		ErrorMsg:	msg,
	}
	c.JSON(http.StatusOK, RestfulAPIResponse{1, err})
}