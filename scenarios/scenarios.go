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
	LeftDecision Decision;
	RightDecision Decision;
}
