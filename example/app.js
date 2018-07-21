'use strict' // warning: code smells

let nameElement = document.getElementById('usr')
let roomElement = document.getElementById('room')
let sendElement = document.getElementById('send')
let chatElement = document.getElementById("chat")
let joinElement = document.getElementById("join")
let msgsElement = document.getElementById('messages')

nameElement.addEventListener("keydown", function(event) {
	if(event.key !== "Enter")
		return
	event.preventDefault()
	join()
})

function join() {
	//let ws = new WebSocket('wss://example.com/chat/room/general')
	let ws = new WebSocket('ws://localhost:8081/general')

	window.addEventListener("beforeunload", function() {
		ws.close()
	})

	ws.onclose = function(e) {
		let reason = 'Unknown reason'
		if(event.code in wsCodeReasons)
			reason = wsCodeReasons[event.code]
		if(event.code == 1010)
			reason += '\nSpecifically, the extensions that are needed are: ' + event.reason
		alert(reason)
	}

	ws.onopen = function() {
		joinElement.style.display = 'none'
		chatElement.style.display = 'block'

		sendElement.onkeyup = function(e) {
			if(e.keyCode == 13) {
				ws.send(JSON.stringify({
					To: roomElement.value,
					Data: {
						Text: sendElement.value,
						From: nameElement.value,
					}
				}))
				e.target.value = ''
			}
		}

		ws.onmessage = function(e) {
			let msg = JSON.parse(e.data)

			let node = document.createElement('p')

			let k = document.createElement('span')
			k.setAttribute('class', 'user')
			k.appendChild(document.createTextNode(msg.Data.From))
			node.appendChild(k)

			let v = document.createElement('span')
			v.setAttribute('class', 'message')
			v.appendChild(document.createTextNode(msg.Data.Text))
			node.appendChild(v)

			prependChild(msgsElement, node)
		}

	}
}

let prependChild = function(c, e) {
	if(c.firstChild)
		c.insertBefore(e, c.firstChild)
	else
		c.appendChild(e)
}

// See http://tools.ietf.org/html/rfc6455#section-7.4.1
let wsCodeReasons = {
	1000: 'Normal closure, meaning that the purpose for which the connection was established has been fulfilled.',
	1001: 'An endpoint is "going away", such as a server going down or a browser having navigated away from a page.',
	1002: 'An endpoint is terminating the connection due to a protocol error',
	1003: 'An endpoint is terminating the connection because it has received a type of data it cannot accept (e.g., an endpoint that understands only text data MAY send this if it receives a binary message).',
	1004: 'Reserved. The specific meaning might be defined in the future.',
	1005: 'No status code was actually present.',
	1006: 'The connection was closed abnormally, e.g., without sending or receiving a Close control frame',
	1007: 'An endpoint is terminating the connection because it has received data within a message that was not consistent with the type of the message (e.g., non-UTF-8 [http://tools.ietf.org/html/rfc3629] data within a text message).',
	1008: 'An endpoint is terminating the connection because it has received a message that "violates its policy". This reason is given either if there is no other sutible reason, or if there is a need to hide specific details about the policy.',
	1009: 'An endpoint is terminating the connection because it has received a message that is too big for it to process.',
	// Note that this status code is not used by the server, because it can fail the WebSocket handshake instead.
	1010: 'An endpoint (client) is terminating the connection because it has expected the server to negotiate one or more extension, but the server didn\'t return them in the response message of the WebSocket handshake.',
	1011: 'A server is terminating the connection because it encountered an unexpected condition that prevented it from fulfilling the request.',
	1015: 'The connection was closed due to a failure to perform a TLS handshake (e.g., the server certificate can\'t be verified).',
}
