package game

import (
	"goobl/scenarios"
)

type Game struct {
	Cross int;
	Population int;
	Sword int;
	Money int;
	Scenario scenario.Scenario;
	ScenarioIndex int;
}
