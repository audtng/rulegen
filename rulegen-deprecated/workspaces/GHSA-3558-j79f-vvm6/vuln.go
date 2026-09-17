package main

		db = global.MustGetGlobalDBByDBName(template.DBName)
	}

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
                    for i := range parts { parts[i] = strings.TrimSpace(parts[i]) }
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
        if v, e := strconv.Atoi(limitStr); e == nil { limitInt = v }
    }
    if offsetStr != "" {
        if v, e := strconv.Atoi(offsetStr); e == nil { offsetInt = v }
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

	var titleKeyMap = make(map[string]string)
	for key, title := range templateInfoMap {
		titleKeyMap[title] = key
	}

	db := global.GVA_DB
	if template.DBName != "" {
		db = global.MustGetGlobalDBByDBName(template.DBName)
	}

	return db.Transaction(func(tx *gorm.DB) error {
		excelTitle := rows[0]
		for i, str := range excelTitle {
			excelTitle[i] = strings.TrimSpace(str)
		}
		values := rows[1:]
		items := make([]map[string]interface{}, 0, len(values))
		for _, row := range values {
			var item = make(map[string]interface{})
			for ii, value := range row {
				if _, ok := titleKeyMap[excelTitle[ii]]; !ok {
					continue // excel中多余的标题，在模板信息中没有对应的字段，因此key为空，必须跳过
				}
				key := titleKeyMap[excelTitle[ii]]
				item[key] = value
			}

			needCreated := tx.Migrator().HasColumn(template.TableName, "created_at")
			needUpdated := tx.Migrator().HasColumn(template.TableName, "updated_at")

			if item["created_at"] == nil && needCreated {
				item["created_at"] = time.Now()
			}
			if item["updated_at"] == nil && needUpdated {
				item["updated_at"] = time.Now()
			}

			items = append(items, item)
		}
		cErr := tx.Table(template.TableName).CreateInBatches(&items, 1000).Error
		return cErr
	})
}

func getColumnName(n int) string {
		mcp.WithDescription(`**GVA代码生成执行器：直接执行代码生成，无需确认步骤**

**核心功能：**
- 根据需求分析和当前的包信息判断是否调用，如果需要调用，则根据入参描述生成json，用于直接生成代码
- 支持批量创建多个模块
- 自动创建包、模块、字典等
- 移除了确认步骤，提高执行效率

**使用场景：**
- 在gva_analyze获取了当前的包信息和字典信息之后，如果已经包含了可以使用的包和模块，那就不要调用本mcp
- 根据分析结果直接生成代码
- 适用于自动化代码生成流程

**批量创建功能：**
- 支持在单个ExecutionPlan中创建多个模块
- modulesInfo字段为数组，可包含多个模块配置
- 一次性处理多个模块的创建和字典生成

**新功能：自动字典创建**
- 当结构体字段使用了字典类型（dictType不为空）时，系统会自动检查字典是否存在
- 如果字典不存在，会自动创建对应的字典及默认的字典详情项
- 字典创建包括：字典主表记录和默认的选项值（选项1、选项2等）

**重要限制：**
- 当needCreatedModules=true时，模块创建会自动生成API和菜单，因此不应再调用api_creator和menu_creator工具
- 只有在单独创建API或菜单（不涉及模块创建）时才使用api_creator和menu_creator工具

重要：ExecutionPlan结构体格式要求（支持批量创建）：
{
  "packageName": "包名(string)",
  "packageType": "package或plugin(string)，如果用户提到了使用插件则创建plugin，如果用户没有特定说明则一律选用package。",
  "needCreatedPackage": "是否需要创建包(bool)",
  "needCreatedModules": "是否需要创建模块(bool)",
  "needCreatedDictionaries": "是否需要创建字典(bool)",
  "packageInfo": {
    "desc": "描述(string)",
    "label": "展示名(string)", 
    "template": "package或plugin(string)，如果用户提到了使用插件则创建plugin，如果用户没有特定说明则一律选用package。",
    "packageName": "包名(string)"
  },
  "modulesInfo": [{
    "package": "包名(string，必然是小写开头)",
    "tableName": "数据库表名(string，使用蛇形命名法)",
    "businessDB": "业务数据库(string)",
    "structName": "结构体名(string)",
    "packageName": "文件名称(string)",
    "description": "中文描述(string)",
    "abbreviation": "简称(string)",
    "humpPackageName": "文件名称 一般是结构体名的小驼峰(string)",
    "gvaModel": "是否使用GVA模型(bool) 固定为true 后续不需要创建ID created_at deleted_at updated_at",
    "autoMigrate": "是否自动迁移(bool)",
    "autoCreateResource": "是否创建资源(bool，默认为false)",
    "autoCreateApiToSql": "是否创建API(bool，默认为true)",
    "autoCreateMenuToSql": "是否创建菜单(bool，默认为true)",
    "autoCreateBtnAuth": "是否创建按钮权限(bool，默认为false)",
    "onlyTemplate": "是否仅模板(bool，默认为false)",
    "isTree": "是否树形结构(bool，默认为false)",
    "treeJson": "树形JSON字段(string)",
    "isAdd": "是否新增(bool) 固定为false",
    "generateWeb": "是否生成前端(bool)",
    "generateServer": "是否生成后端(bool)",
    "fields": [{
      "fieldName": "字段名(string)必须大写开头",
      "fieldDesc": "字段描述(string)",
      "fieldType": "字段类型支持：string（字符串）,richtext（富文本）,int（整型）,bool（布尔值）,float64（浮点型）,time.Time（时间）,enum（枚举）,picture（单图片，字符串）,pictures（多图片，json字符串）,video（视频，字符串）,file（文件，json字符串）,json（JSON）,array（数组）",
      "fieldJson": "JSON标签(string)",
      "dataTypeLong": "数据长度(string)",
      "comment": "注释(string)",
      "columnName": "数据库列名(string)",
      "fieldSearchType": "搜索类型:=/>/</>=/<=/NOT BETWEEN/LIKE/BETWEEN/IN/NOT IN等(string)",
      "fieldSearchHide": "是否隐藏搜索(bool)",
      "dictType": "字典类型(string)",
      "form": "表单显示(bool)",
      "table": "表格显示(bool)",
      "desc": "详情显示(bool)",
      "excel": "导入导出(bool)",
      "require": "是否必填(bool)",
      "defaultValue": "默认值(string)",
      "errorText": "错误提示(string)",
      "clearable": "是否可清空(bool)",
      "sort": "是否排序(bool)",
      "primaryKey": "是否主键(bool)",
      "dataSource": "数据源配置(object) - 用于配置字段的关联表信息，结构：{\"dbName\":\"数据库名\",\"table\":\"关联表名\",\"label\":\"显示字段\",\"value\":\"值字段\",\"association\":1或2(1=一对一,2=一对多),\"hasDeletedAt\":true/false}。\n\n**获取表名提示：**\n- 可在 server/model 和 plugin/xxx/model 目录下查看对应模块的 TableName() 接口实现获取实际表名\n- 例如：SysUser 的表名为 \"sys_users\"，ExaFileUploadAndDownload 的表名为 \"exa_file_upload_and_downloads\"\n- 插件模块示例：Info 的表名为 \"gva_announcements_info\"\n\n**获取数据库名提示：**\n- 主数据库：通常使用 \"gva\"（默认数据库标识）\n- 多数据库：可在 config.yaml 的 db-list 配置中查看可用数据库的 alias-name 字段\n- 如果用户未提及关联多数据库信息 则使用默认数据库 默认数据库的情况下 dbName此处填写为空",
      "checkDataSource": "是否检查数据源(bool) - 启用后会验证关联表的存在性",
      "fieldIndexType": "索引类型(string)"
    }]
  }, {
    "package": "包名(string)",
    "tableName": "第二个模块的表名(string)",
    "structName": "第二个模块的结构体名(string)",
    "description": "第二个模块的描述(string)",
    "...": "更多模块配置..."
  }],
	"dictionariesInfo":[{
		"dictType": "字典类型(string) - 用于标识字典的唯一性",
		"dictName": "字典名称(string) - 必须生成，字典的中文名称",
		"description": "字典描述(string) - 字典的用途说明",
		"status": "字典状态(bool) - true启用，false禁用",
		"fieldDesc": "字段描述(string) - 用于AI理解字段含义并生成合适的选项",
		"options": [{
			"label": "显示名称(string) - 用户看到的选项名",
			"value": "选项值(string) - 实际存储的值",
			"sort": "排序号(int) - 数字越小越靠前"
		}]
	}]
}

注意：
1. needCreatedPackage=true时packageInfo必需
2. needCreatedModules=true时modulesInfo必需
3. needCreatedDictionaries=true时dictionariesInfo必需
4. dictionariesInfo中的options字段可选，如果不提供将根据fieldDesc自动生成默认选项
5. 字典创建会在模块创建之前执行，确保模块字段可以正确引用字典类型
6. packageType只能是"package"或"plugin,如果用户提到了使用插件则创建plugin，如果用户没有特定说明则一律选用package。"
7. 字段类型支持：string（字符串）,richtext（富文本）,int（整型）,bool（布尔值）,float64（浮点型）,time.Time（时间）,enum（枚举）,picture（单图片，字符串）,pictures（多图片，json字符串）,video（视频，字符串）,file（文件，json字符串）,json（JSON）,array（数组）
8. 搜索类型支持：=,!=,>,>=,<,<=,NOT BETWEEN/LIKE/BETWEEN/IN/NOT IN
9. gvaModel=true时自动包含ID,CreatedAt,UpdatedAt,DeletedAt字段
10. **重要**：当gvaModel=false时，必须有一个字段的primaryKey=true，否则会导致PrimaryField为nil错误
11. **重要**：当gvaModel=true时，系统会自动设置ID字段为主键，无需手动设置primaryKey=true
12. 智能字典创建功能：当字段使用字典类型(DictType)时，系统会：
   - 自动检查字典是否存在，如果不存在则创建字典
   - 根据字典类型和字段描述智能生成默认选项，支持状态、性别、类型、等级、优先级、审批、角色、布尔值、订单、颜色、尺寸等常见场景
   - 为无法识别的字典类型提供通用默认选项
13. **模块关联配置**：当需要配置模块间的关联关系时，使用dataSource字段：
   - **dbName**: 关联的数据库名称
   - **table**: 关联的表名
   - **label**: 用于显示的字段名（如name、title等）
   - **value**: 用于存储的值字段名（通常是id）
   - **association**: 关联关系类型（1=一对一关联，2=一对多关联）一对一和一对多的前面的一是当前的实体，如果他只能关联另一个实体的一个，则选用一对一，如果他需要关联多个他的关联实体，则选用一对多。
   - **hasDeletedAt**: 关联表是否有软删除字段
   - **checkDataSource**: 设为true时会验证关联表的存在性
   - 示例：{"dbName":"","table":"sys_users","label":"username","value":"id","association":1,"hasDeletedAt":true}
14. **自动字段类型修正**：系统会自动检查和修正字段类型：
   - 当字段配置了dataSource且association=2（一对多关联）时，系统会自动将fieldType修改为'array'
   - 这确保了一对多关联数据的正确存储和处理
   - 修正操作会记录在日志中，便于开发者了解变更情况`),
        mcp.WithObject("executionPlan",
            mcp.Description("执行计划，包含包信息、模块与字典信息"),
            mcp.Required(),
                },
                "packageType": map[string]interface{}{
                    "type":        "string",
                    "description": "package 或 plugin",
                    "enum":        []string{"package", "plugin"},
                },
                "needCreatedPackage": map[string]interface{}{
                    "type":        "boolean",
                    "description": "是否需要创建包",
                },
                "needCreatedModules": map[string]interface{}{
                    "type":        "boolean",
                    "description": "是否需要创建模块",
                },
                "needCreatedDictionaries": map[string]interface{}{
                    "type":        "boolean",
                    "description": "是否需要创建字典",
                },
                "packageInfo": map[string]interface{}{
                    "type":        "object",
                    "description": "包创建信息",
                    "properties": map[string]interface{}{
                        "desc":        map[string]interface{}{"type": "string", "description": "包描述"},
                        "label":       map[string]interface{}{"type": "string", "description": "展示名"},
                        "template":    map[string]interface{}{"type": "string", "description": "package 或 plugin", "enum": []string{"package", "plugin"}},
                        "packageName": map[string]interface{}{"type": "string", "description": "包名"},
                    },
                },
                "modulesInfo": map[string]interface{}{
                    "type":        "array",
                    "description": "模块配置列表",
                    "items": map[string]interface{}{
                        "type": "object",
                        "properties": map[string]interface{}{
                            "package":            map[string]interface{}{"type": "string", "description": "包名（小写开头）"},
                            "tableName":          map[string]interface{}{"type": "string", "description": "数据库表名（蛇形命名）"},
                            "businessDB":        map[string]interface{}{"type": "string", "description": "业务数据库（可留空表示默认）"},
                            "structName":         map[string]interface{}{"type": "string", "description": "结构体名（大驼峰）"},
                            "packageName":        map[string]interface{}{"type": "string", "description": "文件名称"},
                            "description":        map[string]interface{}{"type": "string", "description": "中文描述"},
                            "abbreviation":       map[string]interface{}{"type": "string", "description": "简称"},
                            "humpPackageName":    map[string]interface{}{"type": "string", "description": "文件名称（小驼峰）"},
                            "gvaModel":           map[string]interface{}{"type": "boolean", "description": "是否使用GVA模型（固定为true）"},
                            "autoMigrate":        map[string]interface{}{"type": "boolean"},
                            "autoCreateResource": map[string]interface{}{"type": "boolean"},
                            "autoCreateApiToSql": map[string]interface{}{"type": "boolean"},
                            "autoCreateMenuToSql": map[string]interface{}{"type": "boolean"},
                            "autoCreateBtnAuth":  map[string]interface{}{"type": "boolean"},
                            "onlyTemplate":       map[string]interface{}{"type": "boolean"},
                            "isTree":             map[string]interface{}{"type": "boolean"},
                            "treeJson":           map[string]interface{}{"type": "string"},
                            "isAdd":              map[string]interface{}{"type": "boolean"},
                            "generateWeb":        map[string]interface{}{"type": "boolean"},
                            "generateServer":     map[string]interface{}{"type": "boolean"},
                            "fields": map[string]interface{}{
                                "type":  "array",
                                "items": map[string]interface{}{
                                    "type": "object",
                                    "properties": map[string]interface{}{
                                        "fieldName":        map[string]interface{}{"type": "string"},
                                        "fieldDesc":        map[string]interface{}{"type": "string"},
                                        "fieldType":        map[string]interface{}{"type": "string"},
                                        "fieldJson":        map[string]interface{}{"type": "string"},
                                        "dataTypeLong":     map[string]interface{}{"type": "string"},
                                        "comment":          map[string]interface{}{"type": "string"},
                                        "columnName":       map[string]interface{}{"type": "string"},
                                        "fieldSearchType":  map[string]interface{}{"type": "string"},
                                        "fieldSearchHide":  map[string]interface{}{"type": "boolean"},
                                        "dictType":         map[string]interface{}{"type": "string"},
                                        "form":             map[string]interface{}{"type": "boolean"},
                                        "table":            map[string]interface{}{"type": "boolean"},
                                        "desc":             map[string]interface{}{"type": "boolean"},
                                        "excel":            map[string]interface{}{"type": "boolean"},
                                        "require":          map[string]interface{}{"type": "boolean"},
                                        "defaultValue":     map[string]interface{}{"type": "string"},
                                        "errorText":        map[string]interface{}{"type": "string"},
                                        "clearable":        map[string]interface{}{"type": "boolean"},
                                        "sort":             map[string]interface{}{"type": "boolean"},
                                        "primaryKey":       map[string]interface{}{"type": "boolean"},
                                        "dataSource": map[string]interface{}{
                                            "type":       "object",
                                            "properties": map[string]interface{}{
                                                "dbName":        map[string]interface{}{"type": "string"},
                                                "table":         map[string]interface{}{"type": "string"},
                                                "label":         map[string]interface{}{"type": "string"},
                                                "value":         map[string]interface{}{"type": "string"},
                                                "association":   map[string]interface{}{"type": "integer"},
                                                "hasDeletedAt":  map[string]interface{}{"type": "boolean"},
                                            },
                                        },
                                        "checkDataSource":   map[string]interface{}{"type": "boolean"},
                                        "fieldIndexType":    map[string]interface{}{"type": "string"},
                                    },
                                },
                            },
                },
                "dictionariesInfo": map[string]interface{}{
                    "type":        "array",
                    "description": "字典创建信息",
                    "items": map[string]interface{}{
                        "type": "object",
                        "properties": map[string]interface{}{
                            "dictType":    map[string]interface{}{"type": "string"},
                            "dictName":    map[string]interface{}{"type": "string"},
                            "description": map[string]interface{}{"type": "string"},
                            "status":      map[string]interface{}{"type": "boolean"},
                            "fieldDesc":   map[string]interface{}{"type": "string"},
                            "options": map[string]interface{}{
                                "type":  "array",
                                "items": map[string]interface{}{
                                    "type": "object",
                                    "properties": map[string]interface{}{
                                        "label": map[string]interface{}{"type": "string"},
                                        "value": map[string]interface{}{"type": "string"},
                                        "sort":  map[string]interface{}{"type": "integer"},
                                    },
                                },
                            },
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
        // 避免重复记录 panic 恢复日志，panic 由 GinRecovery 单独捕捉入库
        if strings.Contains(entry.Message, "[Recovery from panic]") {
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
