package scenario

type Decision struct {
	Magic int;
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
var DeathScenarios []Scenario

func init() {
	Scenarios = []Scenario{
		{
			Prompt: "What spell should we learn first?",
			Image: "MagicGoobl",
			Decisions: 
			[2]Decision{
				{0, 0, 30, 0, "Fireballl for the big boom!"},
				{0, 30, 0, 0, "Create food, so we can have some snacks!"}},
		},
		{
			Prompt: "Would you like to buy some wares?",
			Image: "MerchantGoobl",
			Decisions: 
			[2]Decision{
				{20, 0, 0, -20,"These shells look neat! Could be nice for a ritual"},
				{0, 10, 0, 10,"I'll sell this cookie I have, if you'd like"}},
		},
		{
			Prompt: "There are beetles attacking us!",
			Image: "KnightGoobl",
			Decisions: 
			[2]Decision{
				{0, 0, -30, 0,"Smash them, Goobl!"},
				{0, -30, 0, 0,"RUN AWAYYYYYYY!"}},
		},
	}
	DeathScenarios = []Scenario{
		{
			Prompt: "This story got too booooooring",
			Image: "MagicGoobl",
			Decisions:
			[2]Decision{
				{0, 0, 0, 0, "What....?"},
				{0, 0, 0, 0, "What....?"}},
		},
		{
			Prompt: "HEHEHEHEHE I BECAME GOD",
			Image: "MagicGoobl",
			Decisions:
			[2]Decision{
				{0, 0, 0, 0, "What....?"},
				{0, 0, 0, 0, "What....?"}},
		},
		{
			Prompt: "You haven't been very nice in this story",
			Image: "BaseGoobl",
			Decisions:
			[2]Decision{
				{0, 0, 0, 0, "What....?"},
				{0, 0, 0, 0, "What....?"}},
		},
		{
			Prompt: "Things are tooo crowded now",
			Image: "BaseGoobl",
			Decisions:
			[2]Decision{
				{0, 0, 0, 0, "What....?"},
				{0, 0, 0, 0, "What....?"}},
		},
		{
			Prompt: "They aren't strong enough in this story",
			Image: "KnightGoobl",
			Decisions:
			[2]Decision{
				{0, 0, 0, 0, "What....?"},
				{0, 0, 0, 0, "What....?"}},
		},
		{
			Prompt: "There's too much fighting for me",
			Image: "KnightGoobl",
			Decisions:
			[2]Decision{
				{0, 0, 0, 0, "What....?"},
				{0, 0, 0, 0, "What....?"}},
		},
		{
			Prompt: "The characters should've had more cool things",
			Image: "MerchantGoobl",
			Decisions:
			[2]Decision{
				{0, 0, 0, 0, "What....?"},
				{0, 0, 0, 0, "What....?"}},
		},
		{
			Prompt: "Hehehehe sooooooooo many shiny things",
			Image: "MerchantGoobl",
			Decisions:
			[2]Decision{
				{0, 0, 0, 0, "What....?"},
				{0, 0, 0, 0, "What....?"}},
		},
	}
}
