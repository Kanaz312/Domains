// Modified from https://www.kirupa.com/html5/drag.htm

var dragItem = document.querySelector("#top-card");
var cardBack = document.querySelector("#top-card-back");
var container = document.querySelector("#top-card-container");
var leftAnswer = document.querySelector("#answer-left");
var rightAnswer = document.querySelector("#answer-right");
var active = false;

container.addEventListener("touchstart", dragStart, false);
container.addEventListener("touchend", dragEnd, false);
container.addEventListener("touchmove", drag, false);

container.addEventListener("mousedown", dragStart, false);
container.addEventListener("mouseup", dragEnd, false);
container.addEventListener("mousemove", drag, false);

setTimeout(() => {
	dragItem.style.transform = "rotateY(0deg)";
	dragItem.style.transition = "0.5s ease-in-out";
	cardBack.style.transform = "rotateY(180deg)";
	cardBack.style.transition = "0.5s ease-in-out";
}, 100)

// Returns -1 for left, 0 for neutral, 1 for right
function getSide(x) {
	let w = window.innerWidth
	if (x < w / 4.0) {
		return -1;
	} else if (x > 3.0 * w / 4.0) {
		return 1;
	} else {
		return 0;
	}
}

function dragStart(e) {
	if (e.target === dragItem) {
		active = true;
	}

	dragItem.style.transition = "";
	cardBack.style.transition = "";
}

function dragEnd(e) {
	dragItem.style.transition = "0.4s ease-out";
	active = false;

	let sentRequest = false;

	let x = 0.0;
	let y = 0.0;
	if (e.type === "touchend") {
		x = e.touches[0].clientX;
		y = e.touches[0].clientY;
	} else {
		x = e.clientX;
		y = e.clientY;
	}

	let side = getSide(x);

	let boundingRect = container.getBoundingClientRect()
	let finalX = 0.0;
	let finalY = y - boundingRect.y - (boundingRect.height / 2);

	if (side !== 0)
	{
		fetch("http://localhost:8080/decide", {method : "post", body : JSON.stringify({decision : side})});
		htmx.trigger("#scenario", "game-state-update");
		leftAnswer.removeAttribute("style");
		rightAnswer.removeAttribute("style");
	}

	if (side === -1) {
		finalX = -boundingRect.x - boundingRect.height;
		setTranslate(finalX, finalY, dragItem);
	} else if (side === 1) {
		finalX = window.innerWidth - boundingRect.x + boundingRect.height;
		setTranslate(finalX, finalY, dragItem);
	} else if (side === 0) {
		setTranslate(0.0, 0.0, dragItem);
	}

	setTimeout(() => {
		dragItem.style.transition = "";
		cardBack.style.transition = "";
	},
		400);
}

function drag(e) {
	if (active) {
		e.preventDefault();
		let x;
		let y;
		if (e.type === "touchmove") {
			x = e.touches[0].clientX;
			y = e.touches[0].clientY;
		} else {
			x = e.clientX;
			y = e.clientY;
		}

		let side = getSide(x);
		if (side === -1) {
			leftAnswer.style.display = "block";
			rightAnswer.style.display = "none";
		} else if (side === 1) {
			leftAnswer.style.display = "none";
			rightAnswer.style.display = "block";
		} else {
			leftAnswer.style.display = "none";
			rightAnswer.style.display = "none";
		}

		let boundingRect = container.getBoundingClientRect()
		x -= boundingRect.x + (boundingRect.width / 2);
		y -= boundingRect.y + (boundingRect.height / 2);

		setTranslate(x, y, dragItem);
	}
}

function setTranslate(xPos, yPos, el) {
	el.style.transform = "translate(" + xPos + "px, " + yPos + "px)";
}
