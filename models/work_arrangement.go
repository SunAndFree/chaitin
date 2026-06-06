package models

const (
	WorkTypeTesting   = "测试"
	WorkTypeDelivery  = "交付"
	WorkTypeAfterSales = "售后"
)

const (
	LocationRemote = "远程"
	LocationOnSite = "现场"
)

const (
	PartnerYes = "是"
	PartnerNo  = "否"
)

const (
	ProgressNotStarted = "未开始"
	ProgressInProgress = "进行中"
	ProgressCompleted  = "已完成"
	ProgressPaused     = "已暂停"
	ProgressCancelled  = "已取消"
)

// WorkArrangement represents a work arrangement record.
// id is the internal primary key (auto-increment), never shown to users.
// project_id is the user-facing project/workorder ID.
type WorkArrangement struct {
	ID        int64   `json:"id"`         // internal PK
	ProjectID int64   `json:"project_id"` // 项目ID/工单ID，用户自行填写
	Date      string  `json:"date"`
	Customer  string  `json:"customer"`
	Project   string  `json:"project"`
	WorkType  string  `json:"work_type"`
	Location  string  `json:"location"`
	Partner   string  `json:"partner"`
	Content   string  `json:"content"`
	Duration  float64 `json:"duration"`
	Progress  string  `json:"progress"`
	Notes     string  `json:"notes"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

// FilterParams represents the filter criteria for querying work arrangements
type FilterParams struct {
	DateFrom string `json:"date_from,omitempty"`
	DateTo   string `json:"date_to,omitempty"`
	Customer string `json:"customer,omitempty"`
	Project  string `json:"project,omitempty"`
	WorkType string `json:"work_type,omitempty"`
	Progress string `json:"progress,omitempty"`
}

func ValidWorkTypes() []string {
	return []string{WorkTypeTesting, WorkTypeDelivery, WorkTypeAfterSales}
}

func ValidLocations() []string {
	return []string{LocationRemote, LocationOnSite}
}

func ValidPartners() []string {
	return []string{PartnerYes, PartnerNo}
}

func ValidProgresses() []string {
	return []string{ProgressNotStarted, ProgressInProgress, ProgressCompleted, ProgressPaused, ProgressCancelled}
}

func IsValidWorkType(wt string) bool {
	for _, v := range ValidWorkTypes() {
		if v == wt {
			return true
		}
	}
	return false
}

func IsValidLocation(loc string) bool {
	for _, v := range ValidLocations() {
		if v == loc {
			return true
		}
	}
	return false
}

func IsValidPartner(p string) bool {
	for _, v := range ValidPartners() {
		if v == p {
			return true
		}
	}
	return false
}

func IsValidProgress(p string) bool {
	for _, v := range ValidProgresses() {
		if v == p {
			return true
		}
	}
	return false
}

// GenerateCopyText generates the copy format string.
func (w *WorkArrangement) GenerateCopyText() string {
	if w.Partner == "是" {
		return "【" + w.Location + "】【" + w.WorkType + "】【生态】" + w.Project + "-" + w.Content
	}
	return "【" + w.Location + "】【" + w.WorkType + "】" + w.Project + "-" + w.Content
}
