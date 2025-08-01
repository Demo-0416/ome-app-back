package v1

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"ome-app-back/services"
)

// FoodRecognitionAPI 处理食物识别相关接口
type FoodRecognitionAPI struct {
	recognitionService *services.FoodRecognitionService
}

// NewFoodRecognitionAPI 创建食物识别API处理实例
func NewFoodRecognitionAPI(recognitionService *services.FoodRecognitionService) *FoodRecognitionAPI {
	return &FoodRecognitionAPI{
		recognitionService: recognitionService,
	}
}

// RecognizeFood 分析上传的食物图片
func (a *FoodRecognitionAPI) RecognizeFood(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		responseError(c, http.StatusUnauthorized, "未授权")
		return
	}

	// 获取会话ID
	sessionID := c.PostForm("session_id")

	// 获取上传的文件
	file, err := c.FormFile("food_image")
	if err != nil {
		responseError(c, http.StatusBadRequest, "获取上传文件失败")
		return
	}

	// 检查文件大小
	if file.Size > 10*1024*1024 { // 限制10MB
		responseError(c, http.StatusBadRequest, "文件大小超过限制")
		return
	}

	// 调用服务处理识别
	result, err := a.recognitionService.RecognizeFood(userID, sessionID, file)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "识别食物失败: "+err.Error())
		return
	}

	responseSuccess(c, result)
}

// GetRecognitionByID 获取食物识别记录详情
func (a *FoodRecognitionAPI) GetRecognitionByID(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		responseError(c, http.StatusUnauthorized, "未授权")
		return
	}

	// 获取识别记录ID
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		responseError(c, http.StatusBadRequest, "无效的ID参数")
		return
	}

	// 获取记录详情
	result, err := a.recognitionService.GetRecognitionByID(id)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "获取识别记录失败")
		return
	}

	responseSuccess(c, result)
}

// GetTodayRecognitions 获取用户今日的食物识别记录
func (a *FoodRecognitionAPI) GetTodayRecognitions(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		responseError(c, http.StatusUnauthorized, "未授权")
		return
	}

	// 获取当日记录
	results, err := a.recognitionService.GetUserTodayRecognitions(userID)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "获取识别记录失败")
		return
	}

	responseSuccess(c, results)
}

// SaveRecognitionToNutrition 保存食物识别结果到用户营养摄入
func (a *FoodRecognitionAPI) SaveRecognitionToNutrition(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		responseError(c, http.StatusUnauthorized, "未授权")
		return
	}

	// 获取识别记录ID
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		responseError(c, http.StatusBadRequest, "无效的ID参数")
		return
	}

	// 保存到营养摄入
	err = a.recognitionService.SaveRecognitionToNutrition(id, userID)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "保存到营养摄入失败: "+err.Error())
		return
	}

	responseSuccess(c, nil)
}

// GetAdoptedRecognitions 获取用户已采用的食物识别记录
func (a *FoodRecognitionAPI) GetAdoptedRecognitions(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		responseError(c, http.StatusUnauthorized, "未授权")
		return
	}

	// 获取查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 限制分页参数
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 获取日期范围参数
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	// 调用服务获取记录
	results, err := a.recognitionService.GetAdoptedRecognitions(userID, page, pageSize, startDate, endDate)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "获取识别记录失败: "+err.Error())
		return
	}

	responseSuccess(c, results)
}
