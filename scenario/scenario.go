package scenario

type Decision struct {
	Cross int;
	Population int;
	Sword int;
	Money int;
	Description string;
}

type Scenario struct {
	Prompt string;
	Image string;
	Decisions [2]Decision;
}


var Scenarios []Scenario


func init() {
	Scenarios = []Scenario{
		{
			Prompt: "We would like to construct a cathedral to spread the Good Word",
			Image: "jack",
			Decisions: 
			[2]Decision{
				{30, 10, 0, -40, "Build the church"},
				{-20, 0, 0, 10, "No good, costs too much"}},
		},
		{
			Prompt: "The people are going hungry",
			Image: "queen",
			Decisions: 
			[2]Decision{
				{-10, 40, 0, -10, "Clear the forests to make new farmland"},
				{0, 30, -30, 0, "Steal from our neighbors"}},
		},
		{
			Prompt: "2",
			Image: "ace",
			Decisions: 
			[2]Decision{
				{0, 0, 0, 0,"#"},
				{0, 0, 0, 0,"#"}},
		},
		{
			Prompt: "3",
			Image: "king",
			Decisions: 
			[2]Decision{
				{0, 0, 0, 0,"#"},
				{0, 0, 0, 0,"#"}},
		},
		{
			Prompt: "4",
			Image: "ace",
			Decisions: 
			[2]Decision{
				{0, 0, 0, 0,"#"},
				{0, 0, 0, 0,"#"}},
		},
	}


}
