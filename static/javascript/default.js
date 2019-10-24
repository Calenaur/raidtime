document.addEventListener('DOMContentLoaded', function() {
	var elems = document.querySelectorAll('.modal');
	var instances = M.Modal.init(elems, {
		startingTop: "20%",
		endingTop: "20%",
	});
});