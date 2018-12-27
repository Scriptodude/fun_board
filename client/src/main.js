function getCookie(cname) {
  var name = cname + "=";
  var decodedCookie = decodeURIComponent(document.cookie);
  var ca = decodedCookie.split(';');
  for(var i = 0; i <ca.length; i++) {
    var c = ca[i];
    while (c.charAt(0) == ' ') {
      c = c.substring(1);
    }
    if (c.indexOf(name) == 0) {
      return c.substring(name.length, c.length);
    }
  }
  return "";
}

function setCookie(cname, cvalue) {
  var d = new Date();
  d.setTime(d.getTime() + 1 * 60 * 1000);
  var expires = "expires="+ d.toUTCString();

  console.log(expires)
  document.cookie = cname + "=" + cvalue + ";" + expires;
}

const f = async () => { 
	let body = "{}";
	let id = getCookie("clientId");

	if(id) {
		console.log("Already had id -> " + id)
		body = `{"clientId": ${id}}`;
	}

	const a = await fetch("/connect", 
		{
			method: "POST",
			headers: {
				"Content-Type": "application/json"
			},
			body: body
		});
	const b = await a.json();
	const newId = b.clientId;

	console.log("Old : " + id + ", new : " + newId);
	if(newId !== parseInt(id)) {
		console.log("Received a new Id")
		setCookie("clientId", id);
	}
}

f()