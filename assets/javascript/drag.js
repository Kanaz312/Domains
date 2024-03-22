// Modified from https://www.kirupa.com/html5/drag.htm

document.addEventListener("DOMContentLoaded", initialize, false);

var dragItem;
var container;
var active = false;

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
	if (e.target === dragItem) {
		active = true;
	}

	dragItem.style.transition = "";
}

function dragEnd(e) {
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
		let currentX;
		let currentY;
		if (e.type === "touchmove") {
			currentX = e.touches[0].clientX;
			currentY = e.touches[0].clientY;
		} else {
			currentX = e.clientX;
			currentY = e.clientY;
		}

		let boundingRect = container.getBoundingClientRect()
		currentX -= boundingRect.x + (boundingRect.width / 2);
		currentY -= boundingRect.y + (boundingRect.height / 2);

		setTranslate(currentX, currentY, dragItem);
	}
}

function setTranslate(xPos, yPos, el) {
	el.style.transform = "translate(" + xPos + "px, " + yPos + "px)";
}
