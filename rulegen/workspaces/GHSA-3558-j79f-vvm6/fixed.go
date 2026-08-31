package main

	}
	response.OkWithMessage("文件变更成功", c)
}

// InitDictionary
// @Tags      AutoCodePlugin
// @Summary   打包插件
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}  "打包插件成功"
// @Router    /autoCode/initDictionary [post]
func (a *AutoCodePluginApi) InitDictionary(c *gin.Context) {
	var dictInfo request.InitDictionary
	err := c.ShouldBindJSON(&dictInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = autoCodePluginService.InitDictionary(dictInfo)
	if err != nil {
		global.GVA_LOG.Error("创建初始化Dictionary失败!", zap.Error(err))
		response.FailWithMessage("创建初始化Dictionary失败"+err.Error(), c)
		return
	}
	response.OkWithMessage("文件变更成功", c)
}
package internal

import (
	"context"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	astutil "github.com/flipped-aurora/gin-vue-admin/server/utils/ast"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/stacktrace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
	"time"
)

type ZapCore struct {
}

func (z *ZapCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	for i := 0; i < len(fields); i++ {
		if fields[i].Key == "business" || fields[i].Key == "folder" || fields[i].Key == "directory" {
			syncer := z.WriteSyncer(fields[i].String)
			z.Core = zapcore.NewCore(global.GVA_CONFIG.Zap.Encoder(), syncer, z.level)
		}
	}
	// 先写入原日志目标
	err := z.Core.Write(entry, fields)

	// 捕捉 Error 及以上级别日志并入库，且可提取 zap.Error(err) 的错误内容
	if entry.Level >= zapcore.ErrorLevel {
		// 避免与 GORM zap 写入互相递归：跳过由 gorm logger writer 触发的日志
		if strings.Contains(entry.Caller.File, "gorm_logger_writer.go") {
			return err
		}

		form := "后端"
		level := entry.Level.String()
		// 生成基础信息
		info := entry.Message

		// 提取 zap.Error(err) 内容
		var errStr string
		for i := 0; i < len(fields); i++ {
			f := fields[i]
			if f.Type == zapcore.ErrorType || f.Key == "error" || f.Key == "err" {
				if f.Interface != nil {
					errStr = fmt.Sprintf("%v", f.Interface)
				} else if f.String != "" {
					errStr = f.String
				}
				break
			}
		}
		if errStr != "" {
			info = fmt.Sprintf("%s | 错误: %s", info, errStr)
		}

		// 附加来源与堆栈信息
		if entry.Caller.File != "" {
			info = fmt.Sprintf("%s \n 源文件:%s:%d", info, entry.Caller.File, entry.Caller.Line)
		}
		stack := entry.Stack
		if stack != "" {
			info = fmt.Sprintf("%s \n 调用栈：%s", info, stack)
			// 解析最终业务调用方，并提取其方法源码
			if frame, ok := stacktrace.FindFinalCaller(stack); ok {
				fnName, fnSrc, sLine, eLine, exErr := astutil.ExtractFuncSourceByPosition(frame.File, frame.Line)
				if exErr == nil {
					info = fmt.Sprintf("%s \n 最终调用方法:%s:%d (%s lines %d-%d)\n----- 产生日志的方法代码如下 -----\n%s", info, frame.File, frame.Line, fnName, sLine, eLine, fnSrc)
				} else {
					info = fmt.Sprintf("%s \n 最终调用方法:%s:%d (%s) | extract_err=%v", info, frame.File, frame.Line, fnName, exErr)
				}
			}
		}

		// 使用后台上下文，避免依赖 gin.Context
		ctx := context.Background()
		_ = service.ServiceGroupApp.SystemServiceGroup.SysErrorService.CreateSysError(ctx, &system.SysError{
			Form:  &form,
			Info:  &info,
			Level: level,
		})
	}
	return err
}

func (z *ZapCore) Sync() error {
	--------------------------------------版权声明--------------------------------------
	** 版权所有方：flipped-aurora开源团队 **
	** 版权持有公司：北京翻转极光科技有限责任公司 **
	** 剔除授权标识需购买商用授权：https://plugin.gin-vue-admin.com/license **
	** 感谢您对Gin-Vue-Admin的支持与关注 合法授权使用更有利于项目的长久发展**
`, global.Version, address, address, global.GVA_CONFIG.MCP.SSEPath, address, global.GVA_CONFIG.MCP.MessagePath)
	initServer(address, Router, 10*time.Minute, 10*time.Minute)
// 目前只有Version正式使用 其余为预留
const (
	// Version 当前版本号
	Version = "v2.8.8"
	// AppName 应用名称
	AppName = "Gin-Vue-Admin"
	// Description 应用描述
// @Tag.Description 用户

// @title                       Gin-Vue-Admin Swagger API接口文档
// @version                     v2.8.8
// @description                 使用gin+vue进行极速开发的全栈开发基础平台
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
	)
}

// Handle 处理工具调用
func (d *DictionaryOptionsGenerator) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 解析请求参数
		mcp.WithDescription(`**GVA代码生成执行器：直接执行代码生成，无需确认步骤**

**核心功能：**
根据需求分析和当前的包信息判断是否调用，直接生成代码。支持批量创建多个模块、自动创建包、模块、字典等。

**使用场景：**
在gva_analyze获取了当前的包信息和字典信息之后，如果已经包含了可以使用的包和模块，那就不要调用本mcp。根据分析结果直接生成代码，适用于自动化代码生成流程。

**重要提示：**
- 当needCreatedModules=true时，模块创建会自动生成API和菜单，不应再调用api_creator和menu_creator工具
- 字段使用字典类型时，系统会自动检查并创建字典
- 字典创建会在模块创建之前执行
- 当字段配置了dataSource且association=2（一对多关联）时，系统会自动将fieldType修改为'array'`),
        mcp.WithObject("executionPlan",
            mcp.Description("执行计划，包含包信息、模块与字典信息"),
            mcp.Required(),
                },
                "packageType": map[string]interface{}{
                    "type":        "string",
                    "description": "package 或 plugin，如果用户提到了使用插件则创建plugin，如果用户没有特定说明则一律选用package",
                    "enum":        []string{"package", "plugin"},
                },
                "needCreatedPackage": map[string]interface{}{
                    "type":        "boolean",
                    "description": "是否需要创建包，为true时packageInfo必需",
                },
                "needCreatedModules": map[string]interface{}{
                    "type":        "boolean",
                    "description": "是否需要创建模块，为true时modulesInfo必需",
                },
                "needCreatedDictionaries": map[string]interface{}{
                    "type":        "boolean",
                    "description": "是否需要创建字典，为true时dictionariesInfo必需",
                },
                "packageInfo": map[string]interface{}{
                    "type":        "object",
                    "description": "包创建信息，当needCreatedPackage=true时必需",
                    "properties": map[string]interface{}{
                        "desc":        map[string]interface{}{"type": "string", "description": "包描述"},
                        "label":       map[string]interface{}{"type": "string", "description": "展示名"},
                        "template":    map[string]interface{}{"type": "string", "description": "package 或 plugin，如果用户提到了使用插件则创建plugin，如果用户没有特定说明则一律选用package", "enum": []string{"package", "plugin"}},
                        "packageName": map[string]interface{}{"type": "string", "description": "包名"},
                    },
                },
                "modulesInfo": map[string]interface{}{
                    "type":        "array",
                    "description": "模块配置列表，支持批量创建多个模块",
                    "items": map[string]interface{}{
                        "type": "object",
                        "properties": map[string]interface{}{
                            "package":            map[string]interface{}{"type": "string", "description": "包名（小写开头）"},
                            "tableName":          map[string]interface{}{"type": "string", "description": "数据库表名（蛇形命名法）"},
                            "businessDB":         map[string]interface{}{"type": "string", "description": "业务数据库（可留空表示默认）"},
                            "structName":         map[string]interface{}{"type": "string", "description": "结构体名（大驼峰）"},
                            "packageName":        map[string]interface{}{"type": "string", "description": "文件名称"},
                            "description":        map[string]interface{}{"type": "string", "description": "中文描述"},
                            "abbreviation":       map[string]interface{}{"type": "string", "description": "简称"},
                            "humpPackageName":    map[string]interface{}{"type": "string", "description": "文件名称（小驼峰），一般是结构体名的小驼峰"},
                            "gvaModel":           map[string]interface{}{"type": "boolean", "description": "是否使用GVA模型（固定为true），自动包含ID、CreatedAt、UpdatedAt、DeletedAt字段"},
                            "autoMigrate":        map[string]interface{}{"type": "boolean", "description": "是否自动迁移数据库"},
                            "autoCreateResource": map[string]interface{}{"type": "boolean", "description": "是否创建资源（默认为false）"},
                            "autoCreateApiToSql": map[string]interface{}{"type": "boolean", "description": "是否创建API（默认为true）"},
                            "autoCreateMenuToSql": map[string]interface{}{"type": "boolean", "description": "是否创建菜单（默认为true）"},
                            "autoCreateBtnAuth":  map[string]interface{}{"type": "boolean", "description": "是否创建按钮权限（默认为false）"},
                            "onlyTemplate":       map[string]interface{}{"type": "boolean", "description": "是否仅模板（默认为false）"},
                            "isTree":             map[string]interface{}{"type": "boolean", "description": "是否树形结构（默认为false）"},
                            "treeJson":           map[string]interface{}{"type": "string", "description": "树形JSON字段"},
                            "isAdd":              map[string]interface{}{"type": "boolean", "description": "是否新增（固定为false）"},
                            "generateWeb":        map[string]interface{}{"type": "boolean", "description": "是否生成前端代码"},
                            "generateServer":     map[string]interface{}{"type": "boolean", "description": "是否生成后端代码"},
                            "fields": map[string]interface{}{
                                "type":        "array",
                                "description": "字段列表",
                                "items": map[string]interface{}{
                                    "type": "object",
                                    "properties": map[string]interface{}{
                                        "fieldName":   map[string]interface{}{"type": "string", "description": "字段名（必须大写开头）"},
                                        "fieldDesc":   map[string]interface{}{"type": "string", "description": "字段描述"},
                                        "fieldType":   map[string]interface{}{"type": "string", "description": "字段类型：string（字符串）、richtext（富文本）、int（整型）、bool（布尔值）、float64（浮点型）、time.Time（时间）、enum（枚举）、picture（单图片）、pictures（多图片）、video（视频）、file（文件）、json（JSON）、array（数组）"},
                                        "fieldJson":   map[string]interface{}{"type": "string", "description": "JSON标签"},
                                        "dataTypeLong": map[string]interface{}{"type": "string", "description": "数据长度"},
                                        "comment":     map[string]interface{}{"type": "string", "description": "注释"},
                                        "columnName":  map[string]interface{}{"type": "string", "description": "数据库列名"},
                                        "fieldSearchType": map[string]interface{}{"type": "string", "description": "搜索类型：=、!=、>、>=、<、<=、LIKE、BETWEEN、IN、NOT IN、NOT BETWEEN"},
                                        "fieldSearchHide": map[string]interface{}{"type": "boolean", "description": "是否隐藏搜索"},
                                        "dictType":        map[string]interface{}{"type": "string", "description": "字典类型，使用字典类型时系统会自动检查并创建字典"},
                                        "form":            map[string]interface{}{"type": "boolean", "description": "表单显示"},
                                        "table":           map[string]interface{}{"type": "boolean", "description": "表格显示"},
                                        "desc":            map[string]interface{}{"type": "boolean", "description": "详情显示"},
                                        "excel":           map[string]interface{}{"type": "boolean", "description": "导入导出"},
                                        "require":         map[string]interface{}{"type": "boolean", "description": "是否必填"},
                                        "defaultValue":    map[string]interface{}{"type": "string", "description": "默认值"},
                                        "errorText":       map[string]interface{}{"type": "string", "description": "错误提示"},
                                        "clearable":       map[string]interface{}{"type": "boolean", "description": "是否可清空"},
                                        "sort":            map[string]interface{}{"type": "boolean", "description": "是否排序"},
                                        "primaryKey":      map[string]interface{}{"type": "boolean", "description": "是否主键（gvaModel=false时必须有一个字段为true）"},
                                        "dataSource": map[string]interface{}{
                                            "type":        "object",
                                            "description": "数据源配置，用于配置字段的关联表信息。获取表名提示：可在 server/model 和 plugin/xxx/model 目录下查看对应模块的 TableName() 接口实现获取实际表名（如 SysUser 的表名为 sys_users）。获取数据库名提示：主数据库通常使用 gva（默认数据库标识），多数据库可在 config.yaml 的 db-list 配置中查看可用数据库的 alias-name 字段，如果用户未提及关联多数据库信息则使用默认数据库，默认数据库的情况下 dbName填写为空",
                                            "properties": map[string]interface{}{
                                                "dbName":       map[string]interface{}{"type": "string", "description": "关联的数据库名称（默认数据库留空）"},
                                                "table":        map[string]interface{}{"type": "string", "description": "关联的表名"},
                                                "label":        map[string]interface{}{"type": "string", "description": "用于显示的字段名（如name、title等）"},
                                                "value":        map[string]interface{}{"type": "string", "description": "用于存储的值字段名（通常是id）"},
                                                "association":  map[string]interface{}{"type": "integer", "description": "关联关系类型：1=一对一关联，2=一对多关联。一对一和一对多的前面的一是当前的实体，如果他只能关联另一个实体的一个则选用一对一，如果他需要关联多个他的关联实体则选用一对多"},
                                                "hasDeletedAt": map[string]interface{}{"type": "boolean", "description": "关联表是否有软删除字段"},
                                            },
                                        },
                                        "checkDataSource": map[string]interface{}{"type": "boolean", "description": "是否检查数据源，启用后会验证关联表的存在性"},
                                        "fieldIndexType":  map[string]interface{}{"type": "string", "description": "索引类型"},
                                    },
                                },
                            },
                },
                "dictionariesInfo": map[string]interface{}{
                    "type":        "array",
                    "description": "字典创建信息，字典创建会在模块创建之前执行",
                    "items": map[string]interface{}{
                        "type": "object",
                        "properties": map[string]interface{}{
                            "dictType":    map[string]interface{}{"type": "string", "description": "字典类型，用于标识字典的唯一性"},
                            "dictName":    map[string]interface{}{"type": "string", "description": "字典名称，必须生成，字典的中文名称"},
                            "description": map[string]interface{}{"type": "string", "description": "字典描述，字典的用途说明"},
                            "status":      map[string]interface{}{"type": "boolean", "description": "字典状态：true启用，false禁用"},
                            "fieldDesc":   map[string]interface{}{"type": "string", "description": "字段描述，用于AI理解字段含义并生成合适的选项"},
                            "options": map[string]interface{}{
                                "type":        "array",
                                "description": "字典选项列表（可选，如果不提供将根据fieldDesc自动生成默认选项）",
                                "items": map[string]interface{}{
                                    "type": "object",
                                    "properties": map[string]interface{}{
                                        "label": map[string]interface{}{"type": "string", "description": "显示名称，用户看到的选项名"},
                                        "value": map[string]interface{}{"type": "string", "description": "选项值，实际存储的值"},
                                        "sort":  map[string]interface{}{"type": "integer", "description": "排序号，数字越小越靠前"},
                                    },
                                },
                            },
package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)
				}

				if stack {
					form := "后端"
					info := fmt.Sprintf("Panic: %v\nRequest: %s\nStack: %s", err, string(httpRequest), string(debug.Stack()))
					level := "error"
					_ = service.ServiceGroupApp.SystemServiceGroup.SysErrorService.CreateSysError(context.Background(), &system.SysError{
						Form:  &form,
						Info:  &info,
						Level: level,
					})
					global.GVA_LOG.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				} else {
					form := "后端"
					info := fmt.Sprintf("Panic: %v\nRequest: %s", err, string(httpRequest))
					level := "error"
					_ = service.ServiceGroupApp.SystemServiceGroup.SysErrorService.CreateSysError(context.Background(), &system.SysError{
						Form:  &form,
						Info:  &info,
						Level: level,
					})
					global.GVA_LOG.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
	APIs     []uint `json:"apis"`
}

type InitDictionary struct {
	PlugName     string `json:"plugName"`
	Dictionaries []uint `json:"dictionaries"`
}

type LLMAutoCode struct {
	Prompt string `json:"prompt" form:"prompt" gorm:"column:prompt;comment:提示语;type:text;"` //提示语
	Mode   string `json:"mode" form:"mode" gorm:"column:mode;comment:模式;type:text;"`        //模式
// 导出模板 结构体  SysExportTemplate
type SysExportTemplate struct {
	global.GVA_MODEL
	DBName       string         `json:"dbName" form:"dbName" gorm:"column:db_name;comment:数据库名称;"`                       //数据库名称
	Name         string         `json:"name" form:"name" gorm:"column:name;comment:模板名称;"`                               //模板名称
	TableName    string         `json:"tableName" form:"tableName" gorm:"column:table_name;comment:表名称;"`                //表名称
	TemplateID   string         `json:"templateID" form:"templateID" gorm:"column:template_id;comment:模板标识;"`            //模板标识
	TemplateInfo string         `json:"templateInfo" form:"templateInfo" gorm:"column:template_info;type:text;"`         //模板信息
	SQL          string         `json:"sql" form:"sql" gorm:"column:sql;type:text;comment:自定义导出SQL;"`                    //自定义导出SQL
	ImportSQL    string         `json:"importSql" form:"importSql" gorm:"column:import_sql;type:text;comment:自定义导入SQL;"` //自定义导入SQL
	Limit        *int           `json:"limit" form:"limit" gorm:"column:limit;comment:导出限制"`
	Order        string         `json:"order" form:"order" gorm:"column:order;comment:排序"`
	Conditions   []Condition    `json:"conditions" form:"conditions" gorm:"foreignKey:TemplateID;references:TemplateID;comment:条件"`
package initialize

import (
	"context"
	model "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/plugin-tool/utils"
)

func Dictionary(ctx context.Context) {
	entities := []model.SysDictionary{}
	utils.RegisterDictionaries(entities...)
}
	// initialize.Viper()
	// 安装插件时候自动注册的api数据请到下方法.Api方法中实现
	initialize.Api(ctx)
	// 安装插件时候自动注册的Menu数据请到下方法.Menu方法中实现
	initialize.Menu(ctx)
	// 安装插件时候自动注册的Dictionary数据请到下方法.Dictionary方法中实现
	initialize.Dictionary(ctx)
	initialize.Gorm(ctx)
	initialize.Router(group)
}
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

func RegisterApis(apis ...system.SysApi) {
	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		for _, api := range apis {
			err := tx.Model(system.SysApi{}).Where("path = ? AND method = ? AND api_group = ? ", api.Path, api.Method, api.ApiGroup).FirstOrCreate(&api).Error
	}
}

func RegisterMenus(menus ...system.SysBaseMenu) {
	parentMenu := menus[0]
	otherMenus := menus[1:]
	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
	}

}

func RegisterDictionaries(dictionaries ...system.SysDictionary) {
	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		for _, dict := range dictionaries {
			details := dict.SysDictionaryDetails
			dict.SysDictionaryDetails = nil
			err := tx.Model(system.SysDictionary{}).Where("type = ?", dict.Type).FirstOrCreate(&dict).Error
			if err != nil {
				zap.L().Error("注册字典失败", zap.Error(err), zap.String("type", dict.Type))
				return err
			}
			for _, detail := range details {
				detail.SysDictionaryID = int(dict.ID)
				err = tx.Model(system.SysDictionaryDetail{}).Where("sys_dictionary_id = ? AND value = ?", dict.ID, detail.Value).FirstOrCreate(&detail).Error
				if err != nil {
					zap.L().Error("注册字典详情失败", zap.Error(err), zap.String("value", detail.Value))
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		zap.L().Error("注册字典失败", zap.Error(err))
	}
}

func Pointer[T any](in T) *T {
	return &in
}
package initialize

import (
	"context"
	model "{{.Module}}/model/system"
	"{{.Module}}/plugin/plugin-tool/utils"
)

func Dictionary(ctx context.Context) {
	entities := []model.SysDictionary{}
	utils.RegisterDictionaries(entities...)
}
// initialize.Api(ctx)
// 安装插件时候自动注册的api数据请到下方法.Menu方法中实现并添加如下方法
// initialize.Menu(ctx)
// 安装插件时候自动注册的api数据请到下方法.Dictionary方法中实现并添加如下方法
// initialize.Dictionary(ctx)
func (p *plugin) Register(group *gin.Engine) {
	ctx := context.Background() 
	initialize.Gorm(ctx)
	}
	{
		publicAutoCodeRouter.POST("llmAuto", autoCodeApi.LLMAuto)
		publicAutoCodeRouter.POST("initMenu", autoCodePluginApi.InitMenu)             // 同步插件菜单
		publicAutoCodeRouter.POST("initAPI", autoCodePluginApi.InitAPI)               // 同步插件API
		publicAutoCodeRouter.POST("initDictionary", autoCodePluginApi.InitDictionary) // 同步插件字典
	}
}
						router := strings.Index(threeDirs[k].Name(), "router")
						hasGorm := strings.Index(threeDirs[k].Name(), "gorm")
						response := strings.Index(threeDirs[k].Name(), "response")
						dictionary := strings.Index(threeDirs[k].Name(), "dictionary")
						if gen != -1 && api != -1 && menu != -1 && viper != -1 && plugin != -1 && config != -1 && router != -1 && hasGorm != -1 && response != -1 && dictionary != -1 {
							return nil, nil, nil, errors.Errorf("[filpath:%s]非法模版文件!", four)
						}
						if api != -1 || menu != -1 || viper != -1 || response != -1 || plugin != -1 || config != -1 || dictionary != -1 {
							creates[four] = filepath.Join(global.GVA_CONFIG.AutoCode.Root, global.GVA_CONFIG.AutoCode.Server, "plugin", entity.PackageName, secondDirs[j].Name(), strings.TrimSuffix(threeDirs[k].Name(), ext))
						}
						if gen != -1 {
	os.WriteFile(apiPath, bf.Bytes(), 0666)
	return nil
}

func (s *autoCodePlugin) InitDictionary(dictInfo request.InitDictionary) (err error) {
	dictPath := filepath.Join(global.GVA_CONFIG.AutoCode.Root, global.GVA_CONFIG.AutoCode.Server, "plugin", dictInfo.PlugName, "initialize", "dictionary.go")
	src, err := os.ReadFile(dictPath)
	if err != nil {
		fmt.Println(err)
	}
	fileSet := token.NewFileSet()
	astFile, err := parser.ParseFile(fileSet, "", src, 0)
	arrayAst := ast.FindArray(astFile, "model", "SysDictionary")
	var dictionaries []system.SysDictionary
	err = global.GVA_DB.Preload("SysDictionaryDetails").Find(&dictionaries, "id in (?)", dictInfo.Dictionaries).Error
	if err != nil {
		return err
	}
	dictExpr := ast.CreateDictionaryStructAst(dictionaries)
	arrayAst.Elts = *dictExpr

	var out []byte
	bf := bytes.NewBuffer(out)
	printer.Fprint(bf, fileSet, astFile)

	os.WriteFile(dictPath, bf.Bytes(), 0666)
	return nil
}
// CreateSysError 创建错误日志记录
// Author [yourname](https://github.com/yourname)
func (sysErrorService *SysErrorService) CreateSysError(ctx context.Context, sysError *system.SysError) (err error) {
	if global.GVA_DB == nil {
		return nil
	}
	err = global.GVA_DB.Create(sysError).Error
	return err
}
		db = global.MustGetGlobalDBByDBName(template.DBName)
	}

	// 如果有自定义SQL，则优先使用自定义SQL
	if template.SQL != "" {
		// 将 url.Values 转换为 map[string]interface{} 以支持 GORM 的命名参数
		sqlParams := make(map[string]interface{})
		for k, v := range paramsValues {
			if len(v) > 0 {
				sqlParams[k] = v[0]
			}
		}

		// 执行原生 SQL，支持 @key 命名参数
		err = db.Raw(template.SQL, sqlParams).Scan(&tableMap).Error
		if err != nil {
			return nil, "", err
		}
	} else {
		if len(template.JoinTemplate) > 0 {
			for _, join := range template.JoinTemplate {
				db = db.Joins(join.JOINS + " " + join.Table + " ON " + join.ON)
			}
		}

		db = db.Select(selects).Table(template.TableName)

		filterDeleted := false

		filterParam := paramsValues.Get("filterDeleted")
		if filterParam == "true" {
			filterDeleted = true
		}

		if filterDeleted {
			// 自动过滤主表的软删除
			db = db.Where(fmt.Sprintf("%s.deleted_at IS NULL", template.TableName))

			// 过滤关联表的软删除(如果有)
			if len(template.JoinTemplate) > 0 {
				for _, join := range template.JoinTemplate {
					// 检查关联表是否有deleted_at字段
					hasDeletedAt := sysExportTemplateService.hasDeletedAtColumn(join.Table)
					if hasDeletedAt {
						db = db.Where(fmt.Sprintf("%s.deleted_at IS NULL", join.Table))
					}
				}
			}
		}

		if len(template.Conditions) > 0 {
			for _, condition := range template.Conditions {
				sql := fmt.Sprintf("%s %s ?", condition.Column, condition.Operator)
				value := paramsValues.Get(condition.From)

				if condition.Operator == "IN" || condition.Operator == "NOT IN" {
					sql = fmt.Sprintf("%s %s (?)", condition.Column, condition.Operator)
				}

				if condition.Operator == "BETWEEN" {
					sql = fmt.Sprintf("%s BETWEEN ? AND ?", condition.Column)
					startValue := paramsValues.Get("start" + condition.From)
					endValue := paramsValues.Get("end" + condition.From)
					if startValue != "" && endValue != "" {
						db = db.Where(sql, startValue, endValue)
					}
					continue
				}

				if value != "" {
					if condition.Operator == "LIKE" {
						value = "%" + value + "%"
					}
					db = db.Where(sql, value)
				}
			}
		}
		// 通过参数传入limit
		limit := paramsValues.Get("limit")
		if limit != "" {
			l, e := strconv.Atoi(limit)
			if e == nil {
				db = db.Limit(l)
			}
		}
		// 模板的默认limit
		if limit == "" && template.Limit != nil && *template.Limit != 0 {
			db = db.Limit(*template.Limit)
		}

		// 通过参数传入offset
		offset := paramsValues.Get("offset")
		if offset != "" {
			o, e := strconv.Atoi(offset)
			if e == nil {
				db = db.Offset(o)
			}
		}

		// 获取当前表的所有字段
		table := template.TableName
		orderColumns, err := db.Migrator().ColumnTypes(table)
		if err != nil {
			return nil, "", err
		}

		// 创建一个 map 来存储字段名
		fields := make(map[string]bool)

		for _, column := range orderColumns {
			fields[column.Name()] = true
		}

		// 通过参数传入order
		order := paramsValues.Get("order")

		if order == "" && template.Order != "" {
			// 如果没有order入参，这里会使用模板的默认排序
			order = template.Order
		}

		if order != "" {
			checkOrderArr := strings.Split(order, " ")
			orderStr := ""
			// 检查请求的排序字段是否在字段列表中
			if _, ok := fields[checkOrderArr[0]]; !ok {
				return nil, "", fmt.Errorf("order by %s is not in the fields", order)
			}
			orderStr = checkOrderArr[0]
			if len(checkOrderArr) > 1 {
				if checkOrderArr[1] != "asc" && checkOrderArr[1] != "desc" {
					return nil, "", fmt.Errorf("order by %s is not secure", order)
				}
				orderStr = orderStr + " " + checkOrderArr[1]
			}
			db = db.Order(orderStr)
		}

		err = db.Debug().Find(&tableMap).Error
		if err != nil {
			return nil, "", err
		}
	}

	var rows [][]string
	rows = append(rows, tableTitle)
	for _, exTable := range tableMap {
// PreviewSQL 预览最终生成的 SQL（不执行查询，仅返回 SQL 字符串）
// Author [piexlmax](https://github.com/piexlmax) & [trae-ai]
func (sysExportTemplateService *SysExportTemplateService) PreviewSQL(templateID string, values url.Values) (sqlPreview string, err error) {
	// 解析 params（与导出逻辑保持一致）
	var params = values.Get("params")
	paramsValues, _ := url.ParseQuery(params)

	// 加载模板
	var template system.SysExportTemplate
	err = global.GVA_DB.Preload("Conditions").Preload("JoinTemplate").First(&template, "template_id = ?", templateID).Error
	if err != nil {
		return "", err
	}

	// 解析模板列
	var templateInfoMap = make(map[string]string)
	columns, err := utils.GetJSONKeys(template.TemplateInfo)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal([]byte(template.TemplateInfo), &templateInfoMap)
	if err != nil {
		return "", err
	}
	var selectKeyFmt []string
	for _, key := range columns {
		selectKeyFmt = append(selectKeyFmt, key)
	}
	selects := strings.Join(selectKeyFmt, ", ")

	// 生成 FROM 与 JOIN 片段
	var sb strings.Builder
	sb.WriteString("SELECT ")
	sb.WriteString(selects)
	sb.WriteString(" FROM ")
	sb.WriteString(template.TableName)

	if len(template.JoinTemplate) > 0 {
		for _, join := range template.JoinTemplate {
			sb.WriteString(" ")
			sb.WriteString(join.JOINS)
			sb.WriteString(" ")
			sb.WriteString(join.Table)
			sb.WriteString(" ON ")
			sb.WriteString(join.ON)
		}
	}

	// WHERE 条件
	var wheres []string

	// 软删除过滤
	filterDeleted := false
	if paramsValues != nil {
		filterParam := paramsValues.Get("filterDeleted")
		if filterParam == "true" {
			filterDeleted = true
		}
	}
	if filterDeleted {
		wheres = append(wheres, fmt.Sprintf("%s.deleted_at IS NULL", template.TableName))
		if len(template.JoinTemplate) > 0 {
			for _, join := range template.JoinTemplate {
				if sysExportTemplateService.hasDeletedAtColumn(join.Table) {
					wheres = append(wheres, fmt.Sprintf("%s.deleted_at IS NULL", join.Table))
				}
			}
		}
	}

	// 模板条件（保留与 ExportExcel 同步的解析规则）
	if len(template.Conditions) > 0 {
		for _, condition := range template.Conditions {
			op := strings.ToUpper(strings.TrimSpace(condition.Operator))
			col := strings.TrimSpace(condition.Column)

			// 预览优先展示传入值，没有则展示占位符
			val := ""
			if paramsValues != nil {
				val = paramsValues.Get(condition.From)
			}

			switch op {
			case "BETWEEN":
				startValue := ""
				endValue := ""
				if paramsValues != nil {
					startValue = paramsValues.Get("start" + condition.From)
					endValue = paramsValues.Get("end" + condition.From)
				}
				if startValue != "" && endValue != "" {
					wheres = append(wheres, fmt.Sprintf("%s BETWEEN '%s' AND '%s'", col, startValue, endValue))
				} else {
					wheres = append(wheres, fmt.Sprintf("%s BETWEEN {start%s} AND {end%s}", col, condition.From, condition.From))
				}
			case "IN", "NOT IN":
				if val != "" {
					// 逗号分隔值做简单展示
					parts := strings.Split(val, ",")
					for i := range parts {
						parts[i] = strings.TrimSpace(parts[i])
					}
					wheres = append(wheres, fmt.Sprintf("%s %s ('%s')", col, op, strings.Join(parts, "','")))
				} else {
					wheres = append(wheres, fmt.Sprintf("%s %s ({%s})", col, op, condition.From))
				}
			case "LIKE":
				if val != "" {
					wheres = append(wheres, fmt.Sprintf("%s LIKE '%%%s%%'", col, val))
				} else {
					wheres = append(wheres, fmt.Sprintf("%s LIKE {%%%s%%}", col, condition.From))
				}
			default:
				if val != "" {
					wheres = append(wheres, fmt.Sprintf("%s %s '%s'", col, op, val))
				} else {
					wheres = append(wheres, fmt.Sprintf("%s %s {%s}", col, op, condition.From))
				}
			}
		}
	}

	if len(wheres) > 0 {
		sb.WriteString(" WHERE ")
		sb.WriteString(strings.Join(wheres, " AND "))
	}

	// 排序
	order := ""
	if paramsValues != nil {
		order = paramsValues.Get("order")
	}
	if order == "" && template.Order != "" {
		order = template.Order
	}
	if order != "" {
		sb.WriteString(" ORDER BY ")
		sb.WriteString(order)
	}

	// limit/offset（如果传入或默认值为0，则不生成）
	limitStr := ""
	offsetStr := ""
	if paramsValues != nil {
		limitStr = paramsValues.Get("limit")
		offsetStr = paramsValues.Get("offset")
	}

	// 处理模板默认limit（仅当非0时）
	if limitStr == "" && template.Limit != nil && *template.Limit != 0 {
		limitStr = strconv.Itoa(*template.Limit)
	}

	// 解析为数值，用于判断是否生成
	limitInt := 0
	offsetInt := 0
	if limitStr != "" {
		if v, e := strconv.Atoi(limitStr); e == nil {
			limitInt = v
		}
	}
	if offsetStr != "" {
		if v, e := strconv.Atoi(offsetStr); e == nil {
			offsetInt = v
		}
	}

	if limitInt > 0 {
		sb.WriteString(" LIMIT ")
		sb.WriteString(strconv.Itoa(limitInt))
		if offsetInt > 0 {
			sb.WriteString(" OFFSET ")
			sb.WriteString(strconv.Itoa(offsetInt))
		}
	} else {
		// 当limit未设置或为0时，仅当offset>0才生成OFFSET
		if offsetInt > 0 {
			sb.WriteString(" OFFSET ")
			sb.WriteString(strconv.Itoa(offsetInt))
		}
	}

	return sb.String(), nil
}

// ExportTemplate 导出Excel模板
		return err
	}

	db := global.GVA_DB
	if template.DBName != "" {
		db = global.MustGetGlobalDBByDBName(template.DBName)
	}

	items, err := sysExportTemplateService.parseExcelToMap(rows, templateInfoMap)
	if err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if template.ImportSQL != "" {
			return sysExportTemplateService.importBySQL(tx, template.ImportSQL, items)
		}
		return sysExportTemplateService.importByGORM(tx, template.TableName, items)
	})
}

func (sysExportTemplateService *SysExportTemplateService) parseExcelToMap(rows [][]string, templateInfoMap map[string]string) ([]map[string]interface{}, error) {
	var titleKeyMap = make(map[string]string)
	for key, title := range templateInfoMap {
		titleKeyMap[title] = key
	}

	excelTitle := rows[0]
	for i, str := range excelTitle {
		excelTitle[i] = strings.TrimSpace(str)
	}
	values := rows[1:]
	items := make([]map[string]interface{}, 0, len(values))
	for _, row := range values {
		var item = make(map[string]interface{})
		for ii, value := range row {
			if ii >= len(excelTitle) {
				continue
			}
			if _, ok := titleKeyMap[excelTitle[ii]]; !ok {
				continue // excel中多余的标题，在模板信息中没有对应的字段，因此key为空，必须跳过
			}
			key := titleKeyMap[excelTitle[ii]]
			item[key] = value
		}
		items = append(items, item)
	}
	return items, nil
}

func (sysExportTemplateService *SysExportTemplateService) importBySQL(tx *gorm.DB, sql string, items []map[string]interface{}) error {
	for _, item := range items {
		if err := tx.Exec(sql, item).Error; err != nil {
			return err
		}
	}
	return nil
}

func (sysExportTemplateService *SysExportTemplateService) importByGORM(tx *gorm.DB, tableName string, items []map[string]interface{}) error {
	needCreated := tx.Migrator().HasColumn(tableName, "created_at")
	needUpdated := tx.Migrator().HasColumn(tableName, "updated_at")

	for _, item := range items {
		if item["created_at"] == nil && needCreated {
			item["created_at"] = time.Now()
		}
		if item["updated_at"] == nil && needUpdated {
			item["updated_at"] = time.Now()
		}
	}
	return tx.Table(tableName).CreateInBatches(&items, 1000).Error
}

func getColumnName(n int) string {
	})
	return exists
}

func CreateDictionaryStructAst(dictionaries []system.SysDictionary) *[]ast.Expr {
	var dictElts []ast.Expr
	for i := range dictionaries {
		statusStr := "true"
		if dictionaries[i].Status != nil && !*dictionaries[i].Status {
			statusStr = "false"
		}

		elts := []ast.Expr{
			&ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "Name"},
				Value: &ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("\"%s\"", dictionaries[i].Name)},
			},
			&ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "Type"},
				Value: &ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("\"%s\"", dictionaries[i].Type)},
			},
			&ast.KeyValueExpr{
				Key: &ast.Ident{Name: "Status"},
				Value: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   &ast.Ident{Name: "utils"},
						Sel: &ast.Ident{Name: "Pointer"},
					},
					Args: []ast.Expr{
						&ast.Ident{Name: statusStr},
					},
				},
			},
			&ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "Desc"},
				Value: &ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("\"%s\"", dictionaries[i].Desc)},
			},
		}

		if len(dictionaries[i].SysDictionaryDetails) > 0 {
			var detailElts []ast.Expr
			for _, detail := range dictionaries[i].SysDictionaryDetails {
				detailStatusStr := "true"
				if detail.Status != nil && !*detail.Status {
					detailStatusStr = "false"
				}

				detailElts = append(detailElts, &ast.CompositeLit{
					Type: &ast.SelectorExpr{
						X:   &ast.Ident{Name: "model"},
						Sel: &ast.Ident{Name: "SysDictionaryDetail"},
					},
					Elts: []ast.Expr{
						&ast.KeyValueExpr{
							Key:   &ast.Ident{Name: "Label"},
							Value: &ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("\"%s\"", detail.Label)},
						},
						&ast.KeyValueExpr{
							Key:   &ast.Ident{Name: "Value"},
							Value: &ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("\"%s\"", detail.Value)},
						},
						&ast.KeyValueExpr{
							Key:   &ast.Ident{Name: "Extend"},
							Value: &ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("\"%s\"", detail.Extend)},
						},
						&ast.KeyValueExpr{
							Key: &ast.Ident{Name: "Status"},
							Value: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   &ast.Ident{Name: "utils"},
									Sel: &ast.Ident{Name: "Pointer"},
								},
								Args: []ast.Expr{
									&ast.Ident{Name: detailStatusStr},
								},
							},
						},
						&ast.KeyValueExpr{
							Key:   &ast.Ident{Name: "Sort"},
							Value: &ast.BasicLit{Kind: token.INT, Value: fmt.Sprintf("%d", detail.Sort)},
						},
					},
				})
			}
			elts = append(elts, &ast.KeyValueExpr{
				Key: &ast.Ident{Name: "SysDictionaryDetails"},
				Value: &ast.CompositeLit{
					Type: &ast.ArrayType{Elt: &ast.SelectorExpr{
						X:   &ast.Ident{Name: "model"},
						Sel: &ast.Ident{Name: "SysDictionaryDetail"},
					}},
					Elts: detailElts,
				},
			})
		}

		dictElts = append(dictElts, &ast.CompositeLit{
			Type: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "model"},
				Sel: &ast.Ident{Name: "SysDictionary"},
			},
			Elts: elts,
		})
	}
	return &dictElts
}
//@return: error, string

func BreakPointContinue(content []byte, fileName string, contentNumber int, contentTotal int, fileMd5 string) (string, error) {
	if strings.Contains(fileName, "..") || strings.Contains(fileMd5, "..") {
		return "", errors.New("文件名或路径不合法")
	}
	path := breakpointDir + fileMd5 + "/"
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
//@return: error, string

func MakeFile(fileName string, FileMd5 string) (string, error) {
	if strings.Contains(fileName, "..") || strings.Contains(FileMd5, "..") {
		return "", errors.New("文件名或路径不合法")
	}
	rd, err := os.ReadDir(breakpointDir + FileMd5)
	if err != nil {
		return finishDir + fileName, err
//@return: error

func RemoveChunk(FileMd5 string) error {
	if strings.Contains(FileMd5, "..") {
		return errors.New("路径不合法")
	}
	err := os.RemoveAll(breakpointDir + FileMd5)
	return err
}
		Bucket: aws.String(global.GVA_CONFIG.AwsS3.Bucket),
		Key:    aws.String(filename),
		Body:   f,
		ContentType: aws.String(file.Header.Get("Content-Type")),
	})
	if err != nil {
		global.GVA_LOG.Error("function uploader.Upload() failed", zap.Any("err", err.Error()))
