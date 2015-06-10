document.addEventListener("DOMContentLoaded", function(event) {
	var rights = document.getElementsByClassName('right');
	if (!rights) return;
	var max = 0;
	for (var i = 0; i < rights.length; i++) {
		var e = rights[i];
		var val = e.offsetLeft + e.offsetWidth;
		if (val > max) {
			max = val;
		}
	}
	var envelope = document.getElementById("envelope");
	if (!envelope) return;
	envelope.style.width = "" + max + "px";
	for (var i = 0; i < rights.length; i++) {
		var e = rights[i];
		e.style.float = "right";
	}
});
