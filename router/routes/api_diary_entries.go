package routes

import (
	"mbvlabs/internal/routing"
)

const ApiDiaryPrefix = "/api/diary"

var ApiDiaryThoughtsCurrentWeek = routing.NewSimpleRoute(
	"/thoughts/current-week",
	"api.diary.thoughts.current_week",
	ApiDiaryPrefix,
)
