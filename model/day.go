package model

type Day string

const (
	Monday    Day = "monday"
	Tuesday   Day = "tuesday"
	Wednesday Day = "wednesday"
	Thursday  Day = "thursday"
	Friday    Day = "friday"
	Saturday  Day = "saturday"
	Sunday    Day = "sunday"
)

func GetDayFromInt(day int) Day {
	dayMap := map[int]Day{
		0: Sunday,
		1: Monday,
		2: Tuesday,
		3: Wednesday,
		4: Thursday,
		5: Friday,
		6: Saturday,
	}
	return dayMap[day]
}
