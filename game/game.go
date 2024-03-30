package game

import (
	"goobl/scenarios"
)

type Game struct {
	Cross int;
	Population int;
	Sword int;
	Money int;
	GameScenario scenario.Scenario;
	ScenarioIndex int;
}
