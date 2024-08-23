package models

import "strconv"

// type User struct {
// 	ID
// 	Name  string `json:"name" gorm:"not null;comment:用户名"`
// 	Email string `json:"email" gorm:"not null;default:"";comment:邮箱"`
// 	Timestamps
// 	SoftDeletes
// }
// User 用户模型
// type User struct {
// 	ID
// 	Name       string
// 	Email      string
// 	Containers []Container `gorm:"foreignKey:UserID"`
// 	SoftDeletes
// }

// Challenge 挑战模型
type ChallengeTemplate struct {
	ID
	Name         string `json:"name" gorm:"unique;not null;comment:挑战名称"`
	TemplateYAML string `json:"template_yaml" gorm:"type:text;not null;comment:模板文件"`
	// ConfigYAML   string `json:"service_yaml" gorm:"not null;comment:服务文件"`
	// Timestamps
	// Containers  []Container `gorm:"foreignKey:ChallengeID"`
	Containers []UserContainer `gorm:"foreignKey:ID"`
}

// Container 容器模型
type UserContainer struct {
	ID
	Name      string `json:"name" gorm:"not null;comment:容器名称"`
	Namespace string `json:"namespace" gorm:"not null;comment:命名空间"`
	// Image    string `json:"image" gorm:"not null;comment:镜像名称"`
	Flag       string `json:"flag" gorm:"not null;comment:flag"`
	Status     string `json:"status" gorm:"not null;comment:容器状态"`
	Url        string `json:"url" gorm:"not null;comment:容器地址"`
	UserID     uint   `json:"user_id" gorm:"not null;comment:用户ID"`
	TemplateID uint   `json:"template_id" gorm:"not null;comment:模板ID"`
	Timestamps
	// SoftDeletes
}

func (userContainer UserContainer) GetUid() string {
	return strconv.Itoa(int(userContainer.ID.ID))
}
