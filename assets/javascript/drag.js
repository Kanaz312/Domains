// Modified from https://www.kirupa.com/html5/drag.htm

window.addEventListener("DOMContentLoaded", (e) => {document.querySelector("#scenario").addEventListener("mouseover", initialize, false);}, false);


var dragItem = null;
var container = null;
var leftAnswer = null;
var rightAnswer = null;
var active = false;

function initialize(e) {

	if ((dragItem === null || container === null)) {
		dragItem = document.querySelector("#top-card");
		container = document.querySelector("#top-card-container");
		leftAnswer = document.querySelector("#answer-left");
		rightAnswer = document.querySelector("#answer-right");

		container.addEventListener("touchstart", dragStart, false);
		container.addEventListener("touchend", dragEnd, false);
		container.addEventListener("touchmove", drag, false);

		container.addEventListener("mousedown", dragStart, false);
		container.addEventListener("mouseup", dragEnd, false);
		container.addEventListener("mousemove", drag, false);
	}
}

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
}

function dragEnd(e) {
	setTranslate(0, 0, dragItem);
	dragItem.style.transition = "0.4s ease-out";
	active = false;

	let sentRequest = false;

	let x = 0.0;
	if (e.type === "touchend") {
		x = e.touches[0].clientX;
	} else {
		x = e.clientX;
	}

	let side = getSide(x);

	if (side == -1) {
		fetch("http://localhost:8080/left", {method : "post"});
		htmx.trigger("#stats", "game-state-update");
		htmx.trigger("#scenario", "game-state-update");
		sentRequest = true;
	}
	else if (side === 1) {
		fetch("http://localhost:8080/right", {method : "post"});
		htmx.trigger("#stats", "game-state-update");
		htmx.trigger("#scenario", "game-state-update");
		sentRequest = true;
	}

	if (sentRequest) {
		dragItem = null;
		container = null;

		leftAnswer.removeAttribute("style");
		rightAnswer.removeAttribute("style");
		leftAnswer = null;
		rightAnswer = null;
	}
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
