package constant

// ProductCategory 商品分类结构
type ProductCategory struct {
	Value int    `json:"value"` // 分类ID
	Label string `json:"label"` // 分类名称
}

// 商品分类列表（唯一数据源，前端直接遍历）
var ProductCategoryList = []ProductCategory{
	{Value: 0, Label: "请选择"},
	{Value: 1, Label: "AI助手"},
	{Value: 2, Label: "机器学习"},
	{Value: 3, Label: "计算机视觉"},
	{Value: 4, Label: "数据分析"},
	{Value: 5, Label: "智能硬件"},
	{Value: 6, Label: "自然语言处理"},
}

// GetProductCategoryName 根据分类ID获取分类名称
func GetProductCategoryName(categoryID int) string {
	for _, category := range ProductCategoryList {
		if category.Value == categoryID {
			return category.Label
		}
	}
	return "未知分类"
}
