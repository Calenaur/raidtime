function processResponse(text, response) {
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

function signup(event, type) {
	let xhr = new XMLHttpRequest();
	xhr.onreadystatechange  = function() {
		console.log(this)
		if (this.readyState != 4) 
			return;

		switch(this.status) {
			case 200:
				processResponse("Successfully signed up", this.response);
				break;
			default:
				M.toast({
					html: `Error could not process request`
				});
				break;
		}
	};
	xhr.responseType = "json"
	xhr.open('GET', '/signup/'+event+'/'+type);
	xhr.send();
}

function cancel_signup(event) {
	let xhr = new XMLHttpRequest();
	xhr.onreadystatechange  = function() {
		console.log(this)
		if (this.readyState != 4) 
			return;

		switch(this.status) {
			case 200:
				processResponse("Successfully cancelled signup", this.response);
				break;
			default:
				M.toast({
					html: `Error could not process request`
				});
				break;
		}
	};
	xhr.responseType = "json"
	xhr.open('GET', '/signup/cancel/'+event+'/-1');
	xhr.send();
}