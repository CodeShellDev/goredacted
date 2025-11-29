function shouldActivate(enabled, disabled) {
	const hostname = window.location.hostname.toLowerCase()

	if (disabled.some((d) => hostname === d || hostname.endsWith("." + d))) {
		return false
	}

	if (enabled.length > 0) {
		return enabled.some((d) => hostname === d || hostname.endsWith("." + d))
	}

	return true
}

chrome.storage.sync.get(["enabledDomains", "disabledDomains"], (result) => {
	if (!shouldActivate(result.enabledDomains, result.disabledDomains)) return

	document.addEventListener("paste", async (e) => {
		e.preventDefault()

		chrome.runtime.sendMessage({ action: "fetchParse" }, (res) => {
			let pasteText = ""

			if (res?.success) {
				pasteText = res.data.output
			} else if (res?.error) {
				console.warn("Fetch failed:", res.error)
			}

			const activeEl = document.activeElement

			if (activeEl.tagName === "INPUT" || activeEl.tagName === "TEXTAREA") {
				const start = activeEl.selectionStart
				const end = activeEl.selectionEnd
				const value = activeEl.value

				activeEl.value = value.slice(0, start) + pasteText + value.slice(end)
				activeEl.selectionStart = activeEl.selectionEnd =
					start + pasteText.length
			} else if (activeEl.isContentEditable) {
				const sel = window.getSelection()
				if (!sel.rangeCount) return

				sel.deleteFromDocument()
				sel.getRangeAt(0).insertNode(document.createTextNode(pasteText))
				sel.collapseToEnd()
			} else {
				console.warn(
					"Paste intercepted but target element is unsupported:",
					activeEl
				)
			}
		})
	})
})
