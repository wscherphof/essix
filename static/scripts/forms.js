function alignRight(event) {
	var envelopes = document.getElementsByClassName("envelope");
	if (!envelopes) return;
	for (var i = 0; i < envelopes.length; i++) {
		var envelope = envelopes[i];
		var rights = envelope.querySelectorAll('.right');
		if (!rights) continue;
		var max = 0;
		for (var i = 0; i < rights.length; i++) {
			var elem = rights[i];
			var val = elem.offsetLeft + elem.offsetWidth;
			if (val > max) {
				max = val;
			}
		}
		for (var i = 0; i < rights.length; i++) {
			var elem = rights[i];
			var offset = max - elem.offsetLeft - elem.offsetWidth;
			elem.style.position = "relative";
			elem.style.left = offset + "px";
		}
	}
}

document.addEventListener("DOMContentLoaded", alignRight);
