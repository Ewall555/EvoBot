package model

type Staff struct {
	BaseModel
	StaffID   string   `gorm:"column:staffid;type:varchar(255);not null;unique" json:"staffid"`
	StaffName string   `gorm:"column:staffname;type:varchar(255);" json:"staffname"`
	Number    string   `gorm:"column:number;type:varchar(255);" json:"number"`
	Email     string   `gorm:"column:email;type:varchar(255);" json:"email"`
	Role      string   `gorm:"column:role;type:varchar(255);" json:"role"` // 售前，销售，售后
	Policies  []Policy `gorm:"many2many:staff_policy;foreignKey:ID;joinForeignKey:StaffID;References:ID;joinReferences:PolicyID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;" json:"staff_policy"`
}
