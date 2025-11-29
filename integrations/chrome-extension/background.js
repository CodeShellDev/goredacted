chrome.runtime.onMessage.addListener((msg, sender, sendResponse) => {
	if (msg.action === "fetchParse") {
		;(async () => {
			try {
				const result = await chrome.storage.sync.get("httpEndpoint")
				const endpoint = result.httpEndpoint
				if (!endpoint) {
					sendResponse({
						success: false,
						error: "HTTP endpoint not configured",
					})
					return
				}

				const resp = await fetch(endpoint, { method: "GET" })
				const data = await resp.json()
				sendResponse({ success: true, data })
			} catch (err) {
				sendResponse({ success: false, error: err.message })
			}
		})()

		return true
	}
})
