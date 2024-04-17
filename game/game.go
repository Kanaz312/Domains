package game

import (
	"goobl/scenario"
	"log"
)

type Game struct {
	Magic int;
	Population int;
	Sword int;
	Money int;
	ScenarioIndex int;
}

func (g *Game) IsDead() bool {
	return g.Magic <= 0 || g.Magic >= 100 ||
		g.Population <= 0 || g.Population >= 100 ||
		g.Sword <= 0 || g.Sword >= 100 ||
		g.Money <= 0 || g.Money >= 100
}

func (g *Game) ApplyDecision(decision *scenario.Decision, scenarios []scenario.Scenario) {
	g.Magic += decision.Magic
	g.Population += decision.Population
	g.Sword += decision.Sword
	g.Money += decision.Money

	if !g.IsDead() {
		g.ScenarioIndex = (g.ScenarioIndex + 1) % len(scenarios)
	}
}

func (g *Game) GetDeathScenario() *scenario.Scenario {
	if  g.Magic <= 0 {
		return &scenario.DeathScenarios[0]
	} else if g.Magic >= 100 {
		return &scenario.DeathScenarios[1]
	} else if g.Population <= 0 {
		return &scenario.DeathScenarios[2]
	} else if g.Population >= 100 {
		return &scenario.DeathScenarios[3]
	} else if g.Sword <= 0 {
		return &scenario.DeathScenarios[4]
	} else if g.Sword >= 100 {
		return &scenario.DeathScenarios[5]
	} else if g.Money <= 0 {
		return &scenario.DeathScenarios[6]
	} else if g.Money >= 100 {
		return &scenario.DeathScenarios[7]
	} else {
		log.Println("FillDeathScenario: Game State Appears non-dead")
		return nil
	}
}
