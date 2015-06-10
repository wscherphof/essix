function max(list) {
	var max = 0;
	if (list) {
		for (var i = 0; i < list.length; i++) {
			var elem = list[i];
			var val = elem.offsetLeft + elem.offsetWidth;
			if (val > max) {
				max = val;
			}
		}
	}
	return max;
}

function position(list, max) {
	if (list) {
		for (var i = 0; i < list.length; i++) {
			var elem = list[i];
			var offset = max - elem.offsetLeft - elem.offsetWidth;
			elem.style.position = "relative";
			elem.style.left = offset + "px";
		}
	}
}

function alignRight(event) {
	var envelopes = document.getElementsByClassName("envelope");
	if (!envelopes) return;
	for (var i = 0; i < envelopes.length; i++) {
		var envelope = envelopes[i];
		var rights = envelope.querySelectorAll('.right');
		var right2s = envelope.querySelectorAll('.right2');
		var maxRights = max(rights);
		var maxRight2s = max(right2s);
		position(rights, maxRights);
		position(right2s, maxRight2s);
	}
}

document.addEventListener("DOMContentLoaded", alignRight);
