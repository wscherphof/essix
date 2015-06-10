function scorePassword(pass) {
    var score = 0;
    if (!pass)
        return score;

    // award every unique letter until 5 repetitions
    var letters = new Object();
    for (var i=0; i<pass.length; i++) {
        letters[pass[i]] = (letters[pass[i]] || 0) + 1;
        score += 5.0 / letters[pass[i]];
    }

    // bonus points for mixing it up
    var variations = {
        digits: /\d/.test(pass),
        lower: /[a-z]/.test(pass),
        upper: /[A-Z]/.test(pass),
        nonWords: /\W/.test(pass),
    }

    variationCount = 0;
    for (var check in variations) {
        variationCount += (variations[check] == true) ? 1 : 0;
    }
    score += (variationCount - 1) * 10;

    return parseInt(score);
}

function judgePassword(value, verdictId) {
    var score = scorePassword(value);
    var color = 'red';
    if (score > 60) {
    	color = 'yellow';
    }
    if (score > 80) {
     	color = 'lime';
    }
    var verdict = document.getElementById(verdictId);
    if (verdict) {
		verdict.innerHTML = '&nbsp;&#8226;&nbsp;';
		verdict.style.backgroundColor = 'silver';
		verdict.style.color = color;
    }
}
