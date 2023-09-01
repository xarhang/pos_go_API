package controller

import (
	"errors"
	"net/http"
	"pos-go/db"
	"pos-go/dto"
	"pos-go/model"
	pubfunc "pos-go/pubFunc"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Permission struct{}

func (p Permission) FindAll(ctx *gin.Context) {
	mpath := ctx.Query("path")

	var permission []dto.VPermissionResponse
	query := db.Conn
	if mpath != "" {
		query = query.Where("path = ?", &mpath)
	}
	query.Find(&permission)
	var result []dto.VPermissionResponse
	for _, pms := range permission {
		result = append(result, dto.VPermissionResponse{
			// ID:         pms.ID,
			// Ptype:      pms.Ptype,
			Rules:      pms.Rules,
			Path:       pms.Path,
			Permission: pms.Permission,
		})
	}
	if result == nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "not found"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}
func (p Permission) FindOne(ctx *gin.Context) {
	rule := ctx.Param("rule")
	var permission []dto.VPermissionResponse
	db.Conn.Find(&permission, "rules=?", rule)
	// sql := `select rules,path,GROUP_CONCAT(permission SEPARATOR ',') permission FROM v_permissions WHERE rules='` +
	// 	rule + `'` +
	// 	`group by rules,path`
	// println(sql)
	// db.Conn.
	// 	Raw(sql).
	// 	Scan(&permission)

	var result []dto.VPermissionResponse
	for _, pms := range permission {
		result = append(result, dto.VPermissionResponse{
			Rules:      pms.Rules,
			Path:       pms.Path,
			Permission: pms.Permission,
		})
	}
	if result == nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "could not be found "})
		return
	}
	ctx.JSON(http.StatusOK, result)

}
func (p Permission) Create(ctx *gin.Context) {
	var permission dto.VPermissionrequest
	if err := ctx.ShouldBindJSON(&permission); err != nil {
		if err.Error() == "EOF" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no request body"})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !pubfunc.IsValidPtype(permission.Ptype) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "incorrect ptype input"})
		return
	}
	if !pubfunc.IsValidPermission(permission.Permission) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "incorrect Permission input"})
		return
	}
	var rule model.Rule
	ruleCheck := db.Conn.First(&rule, "rule_name=?", permission.Rules)
	if errors.Is(ruleCheck.Error, gorm.ErrRecordNotFound) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "incorrect rules input"})
		return
	}
	var vPermission model.VPermission
	doubCheck := db.Conn.Where("ptype='p' and rules=? and path=? and permission=?",
		permission.Rules,
		permission.Path,
		permission.Permission).First(&vPermission)
	if doubCheck.RowsAffected > 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Already exist!"})
		return
	}
	VPermissionVar := model.VPermission{
		Ptype:      permission.Ptype,
		Rules:      permission.Rules,
		Path:       permission.Path,
		Permission: permission.Permission,
	}
	if err := db.Conn.Create(&VPermissionVar).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "success", "permission_data": VPermissionVar})
}

func (p Permission) Delete(ctx *gin.Context) {
	var permission dto.VPermissionrequest
	if err := ctx.ShouldBindJSON(&permission); err != nil {
		if err.Error() == "EOF" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no request body"})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !pubfunc.IsValidPtype(permission.Ptype) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "incorrect ptype input"})
		return
	}
	if !pubfunc.IsValidPermission(permission.Permission) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "incorrect Permission input"})
		return
	}
	var rule model.Rule
	ruleCheck := db.Conn.First(&rule, "rule_name=?", permission.Rules)
	if errors.Is(ruleCheck.Error, gorm.ErrRecordNotFound) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "incorrect rules input"})
		return
	}
	// var vPermission model.VPermission
	// doubCheck := db.Conn.Where("ptype=? and rules=? and path=? and permission=?",
	// 	permission.Ptype,
	// 	permission.Rules,
	// 	permission.Path,
	// 	permission.Permission).First(&vPermission)
	// if doubCheck.RowsAffected == 0 {
	// 	ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "incorrect in put!"})
	// 	return
	// }
	VPermissionVar := model.VPermission{
		Ptype:      permission.Ptype,
		Rules:      permission.Rules,
		Path:       permission.Path,
		Permission: permission.Permission,
	}
	err := db.Conn.Where("ptype=? and rules=? and path=? and permission=?",
		permission.Ptype,
		permission.Rules,
		permission.Path,
		permission.Permission).Delete(&VPermissionVar)
	if err.RowsAffected == 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no record found"})
		return
	}
	ctx.JSON(http.StatusAccepted, gin.H{"message": "success", "permission_data": VPermissionVar})
}
