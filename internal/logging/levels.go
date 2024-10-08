package logging

type LogStoreLevel string

const (
	LogStoreLevelDelete       LogStoreLevel = "delete"
	LogStoreLevelCreate       LogStoreLevel = "create"
	LogStoreLevelUpdate       LogStoreLevel = "update"
	LogStoreLevelInstanceView LogStoreLevel = "instance_view"
	LogStoreLevelListView     LogStoreLevel = "list_view"
	LogStoreLevelPanelView    LogStoreLevel = "panel_view"
)

var levelsHierarchy = map[LogStoreLevel]int{
	LogStoreLevelDelete:       1,
	LogStoreLevelCreate:       2,
	LogStoreLevelUpdate:       3,
	LogStoreLevelInstanceView: 4,
	LogStoreLevelListView:     5,
	LogStoreLevelPanelView:    6,
}

func (l LogStoreLevel) AssessLevel(assessmentLevel LogStoreLevel) bool {
	currentLevelRank, currentExists := levelsHierarchy[l]
	assessmentLevelRank, assessmentExists := levelsHierarchy[assessmentLevel]

	if !currentExists || !assessmentExists {
		return true
	}

	return currentLevelRank >= assessmentLevelRank
}
