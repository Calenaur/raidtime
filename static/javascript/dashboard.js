var calendar
var monthOffset = 0
var calendarXHR = null
var me = null
var allEvents = {}
var monthNames = [
	"January", "February", "March",
	"April", "May", "June", "July",
	"August", "September", "October",
	"November", "December"
];

function handle_response(text, response) {
	if (response.code == 200) {
		M.toast({
			html: text
		});
		return;
	}
	M.toast({
		html: response.error
	});
}

function signup(type) {
	modal = document.getElementById("event-modal");
	event = allEvents[parseInt(modal.getAttribute("event-id"))];
	if (hasSignedEvent(me, event))
		return;

	let xhr = new XMLHttpRequest();
	xhr.onreadystatechange  = function() {
		if (this.readyState != 4) 
			return;

		switch(this.status) {
			case 200:
				handle_response("Successfully signed up", this.response);
				if (this.response.signup != null)
					event.signups.push(this.response.signup);
				break;
			default:
				M.toast({
					html: `Error could not process request`
				});
				break;
		}
	};
	xhr.responseType = "json"
	xhr.open('GET', '/signup/'+event.id+'/'+type);
	xhr.send();
}

function cancel_signup() {
	modal = document.getElementById("event-modal");
	event = allEvents[parseInt(modal.getAttribute("event-id"))];
	if (!hasSignedEvent(me, event))
		return;

	let xhr = new XMLHttpRequest();
	xhr.onreadystatechange  = function() {
		if (this.readyState != 4) 
			return;

		switch(this.status) {
			case 200:
				handle_response("Successfully cancelled signup", this.response);
				for (var i=0; i<event.signups.length; i++) {
					s = event.signups[i];
					if (s.user.id == me)
						event.signups.splice(i, 1);
				}
				break;
			default:
				M.toast({
					html: `Error could not process request`
				});
				break;
		}
	};
	xhr.responseType = "json"
	xhr.open('GET', '/signup/cancel/'+event.id+'/-1');
	xhr.send();
}

function getDateByOffset(offset=0) {
	now = new Date();
	date = new Date(now.setMonth(now.getMonth()+offset));
	return date
}

function getDaysInMonth(date) {
	return new Date(date.getFullYear(), date.getMonth()+1, 0).getDate()
}

function getFirstDayOfMonth(date) {
	return new Date(date.getFullYear(), date.getMonth(), 1).getDay()
}

function hasSignedEvent(id, event) {
	for (signee of event.signups)
		if (signee.user.id == id)
			return true;
	return false;
}

function init() {
	calendar = document.getElementById("calendar-body").children
	url = location.href;
	monthOffset = parseInt(url.substring(url.indexOf("#")+1));
	if (isNaN(monthOffset))
		monthOffset = 0;
}

function load_calendar(offset) {
	display_calendar(offset);
	if (calendarXHR != null) {
		calendarXHR.abort();
	}
	calendarXHR = new XMLHttpRequest();
	calendarXHR.onreadystatechange  = function() {
		if (this.readyState != 4) 
			return;

		if (this.status == 200)
			if (this.response.code == 200) {
				display_calendar(offset, this.response);
				me = this.response.me;
				return;
			}
	};
	calendarXHR.responseType = "json"
	calendarXHR.open('GET', '/calendar/'+offset);
	calendarXHR.send();
}

function display_calendar(offset, data=null) {
	date = getDateByOffset(offset);
	daysInMonth = getDaysInMonth(date);
	firstDay = getFirstDayOfMonth(date);
	weeks = Math.ceil((daysInMonth+firstDay)/7);
	document.getElementById("calendar-title").innerHTML = monthNames[date.getMonth()] + " " + date.getFullYear();
	eventMap = {}
	if (data != null) {
		if (data.events != null)
			for (event of data.events) {
				event.date = new Date(event.date);
				k = event.date.getDate();
				if (eventMap[k] == null)
					eventMap[k] = [];
				eventMap[k].push(event);
			}
	} 

	weekI = 0;
	dayI = 0;
	dayN = 0;
	for (week of calendar) {
		if (weekI < weeks) {
			week.classList.remove("hidden-week");
			for (day of week.children) {
				dayName = day.getElementsByClassName("day")[0];
				events = day.getElementsByClassName("events")[0];
				if (dayI >= firstDay && dayN < daysInMonth) {
					dayN++;
					dayName.innerHTML = dayN;
					while (events.firstChild)
						events.removeChild(events.firstChild);

					if (eventMap[dayN] != null) {
						for (event of eventMap[dayN]) {
							allEvents[event.id] = event
							events.appendChild(create_event_element(event));
						}
					}
				} else {
					dayName.innerHTML = "";
				}
				dayI++
			}
		} else {
			week.classList.add("hidden-week");
		}
		weekI++
	}
}

function create_event_element(event) {
	el = document.createElement("div");
	el.classList.add("event");
	el.setAttribute("event-id", event.id);
	el.onclick = open_event
	el.style.backgroundColor = event.tag.color;
	el.innerHTML = event.name;
	return el;
}

function create_signup_element(signee) {
	el = document.createElement("li");

	reasonEl = document.createElement("span");
	reasonEl.innerHTML = signee.reason.description + " - ";

	nameEl = document.createElement("span");
	nameEl.style.color = signee.user.class.color;
	nameEl.innerHTML = signee.user.username;

	rankEl = document.createElement("span");
	rankEl.innerHTML = "&lt" + signee.user.guild_rank.name + "&gt";

	el.appendChild(reasonEl);
	el.appendChild(nameEl);
	el.appendChild(rankEl);

	return el;
}

function open_event() {
	modal = document.getElementById("event-modal");
	id = parseInt(this.getAttribute("event-id"));
	modal.setAttribute("event-id", id);
	event = allEvents[id];
	console.log(event);

	if (event == null)
		return;

	eventNameEl = modal.getElementsByClassName("event-name")[0];
	signeesEl = modal.getElementsByClassName("event-signees")[0];
	creatorEl = modal.getElementsByClassName("event-creator")[0];
	labelEl = modal.getElementsByClassName("event-label")[0];
	signupEl = modal.getElementsByClassName("event-signup")[0];
	cacelSignupEl = modal.getElementsByClassName("event-cancel-signup")[0];

	eventNameEl.innerHTML = event.name;
	creatorEl.innerHTML = event.creator.username;
	labelEl.style.color = event.tag.color;
	if (hasSignedEvent(me, event)) {
		signupEl.style.display = "none";
		cacelSignupEl.style.display = "inline";
	} else {
		signupEl.style.display = "inline";
		cacelSignupEl.style.display = "none";
	}

	event.signups.sort(signee_sort);
	while (signeesEl.firstChild)
		signeesEl.removeChild(signeesEl.firstChild);

	for (signee of event.signups)
		signeesEl.appendChild(create_signup_element(signee));

	M.Modal.init(modal, null).open();
}

function signee_sort(s1, s2) {
	if (s1.reason.will_attend && !s2.reason.will_attend)
		return 1;

	if (!s1.reason.will_attend && s2.reason.will_attend)
		return -1;

	if (s1.reason.description > s2.reason.description)
		return 1;

	if (s1.reason.description < s2.reason.description)
		return -1;

	if (s1.user.class.name > s2.user.class.name)
		return 1;

	if (s1.user.class.name < s2.user.class.name)
		return -1;

	if (s1.user.guild_rank.name > s2.user.guild_rank.name)
		return 1;

	if (s1.user.guild_rank.name < s2.user.guild_rank.name)
		return -1;

	if (s1.user.name > s2.user.name)
		return 1;

	if (s1.user.name < s2.user.name)
		return -1;

	return 0;
}

function next_month() {
	var url = location.href;
    location.href = "#"+(++monthOffset);
	load_calendar(monthOffset)
}

function prev_month() {
	var url = location.href;
    location.href = "#"+(--monthOffset);
	load_calendar(monthOffset)
}

document.addEventListener('DOMContentLoaded', function() {
	init();
	load_calendar(monthOffset);
});