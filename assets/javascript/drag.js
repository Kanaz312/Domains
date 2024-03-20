// Modified from https://www.kirupa.com/html5/drag.htm

document.addEventListener("DOMContentLoaded", initialize, false);

// const gameStateUpdateEvent = new Event("game-state-update");

var dragItem;
var container;
var active = false;
var currentX;
var currentY;
var initialX;
var initialY;

function initialize(e) {

	dragItem = document.querySelector("#top-card");
	container = document.querySelector("#top-card-container");

	container.addEventListener("touchstart", dragStart, false);
	container.addEventListener("touchend", dragEnd, false);
	container.addEventListener("touchmove", drag, false);

	container.addEventListener("mousedown", dragStart, false);
	container.addEventListener("mouseup", dragEnd, false);
	container.addEventListener("mousemove", drag, false);
}

function dragStart(e) {
	if (e.type === "touchstart") {
		initialX = e.touches[0].clientX;
		initialY = e.touches[0].clientY;
	} else {
		initialX = e.clientX;
		initialY = e.clientY;
	}

	if (e.target === dragItem) {
		active = true;
	}

	dragItem.style.transition = "";
}

function dragEnd(e) {
	initialX = 0;
	initialY = 0;


	setTranslate(0, 0, dragItem);
	dragItem.style.transition = "0.4s ease-out";
	active = false;

	let w = window.innerWidth

	if (e.clientX < w / 4.0)
	{
		fetch("http://localhost:8080/left", {method : "post"});
	}
	else if (e.clientX > 3.0 * w / 4.0)
	{
		fetch("http://localhost:8080/right", {method : "post"});
	}

	htmx.trigger("#stats", "game-state-update");
}

function drag(e) {
	if (active) {

		e.preventDefault();

		if (e.type === "touchmove") {
			currentX = e.touches[0].clientX - initialX;
			currentY = e.touches[0].clientY - initialY;
		} else {
			currentX = e.clientX - initialX;
			currentY = e.clientY - initialY;
		}

		setTranslate(currentX, currentY, dragItem);
	}
}

function setTranslate(xPos, yPos, el) {
	el.style.transform = "translate(" + xPos + "px, " + yPos + "px)";
}
