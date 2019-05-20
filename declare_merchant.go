package wego

// MerchantCategoryInfo ...
type MerchantCategoryInfo struct {
	Errcode int64  `json:"errcode"`
	Errmsg  string `json:"errmsg"`
	Data    Data   `json:"data"`
}

// Data ...
type Data struct {
	AllCategoryInfo AllCategoryInfo `json:"all_category_info"`
}

// AllCategoryInfo ...
type AllCategoryInfo struct {
	Categories []Category `json:"categories"`
}

// Category ...
type Category struct {
	ID            int64    `json:"id"`
	Name          string   `json:"name"`
	Level         int64    `json:"level"`
	Children      []int64  `json:"children"`
	Father        *int64   `json:"father,omitempty"`
	Qualify       *Qualify `json:"qualify,omitempty"`
	Scene         *int64   `json:"scene,omitempty"`
	SensitiveType *int64   `json:"sensitive_type,omitempty"`
}

// Qualify ...
type Qualify struct {
	ExterList []ExterList `json:"exter_list"`
}

// ExterList ...
type ExterList struct {
	InnerList []InnerList `json:"inner_list"`
}

// InnerList ...
type InnerList struct {
	Name string `json:"name"`
}

// ApplyMerchantInfo ...
type ApplyMerchantInfo struct {
	FirstCatID        int64  `json:"first_catid"`
	SecondCatID       int64  `json:"second_catid"`
	QualificationList string `json:"qualification_list"`
	HeadimgMediaID    string `json:"headimg_mediaid"`
	Nickname          string `json:"nickname"`
	Intro             string `json:"intro"`
	OrgCode           string `json:"org_code"`
	OtherFiles        string `json:"other_files,omitempty"`
}
