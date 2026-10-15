;(() => {
	console.info("[Goredacted] Content script loaded:", location.href)

	if (!globalThis.GoredactedDomainRules) {
		console.error("[Goredacted] Shared domain rules failed to load")
		return
	}

	const { DEFAULTS, isDomainEnabled } = GoredactedDomainRules

	let settings = null

	async function loadSettings() {
		settings = await chrome.storage.local.get(DEFAULTS)

		console.info("[Goredacted] Settings loaded:", settings)
		console.info("[Goredacted] Hostname:", location.hostname)
		console.info(
			"[Goredacted] Domain enabled:",
			isDomainEnabled(location.hostname, settings),
		)
	}

	function isEnabledHere() {
		if (!settings) {
			console.warn("[Goredacted] Settings have not loaded")
			return false
		}

		const enabled = isDomainEnabled(location.hostname, settings)
		console.info("[Goredacted] isEnabledHere:", enabled)
		return enabled
	}

	function getEditableTarget(node) {
		if (!(node instanceof Element)) return null

		const editable = node.closest(
			'textarea, input:not([type]), input[type="text"], input[type="search"], ' +
				'input[type="email"], input[type="url"], input[type="tel"], [contenteditable="true"], [contenteditable="plaintext-only"]',
		)

		return editable instanceof HTMLElement ? editable : null
	}

	function saveSelection(target) {
		if (
			target instanceof HTMLTextAreaElement ||
			target instanceof HTMLInputElement
		) {
			return {
				kind: "text-control",
				start: target.selectionStart ?? target.value.length,
				end: target.selectionEnd ?? target.value.length,
				direction: target.selectionDirection || "none",
			}
		}

		if (target?.isContentEditable) {
			const selection = window.getSelection()

			if (selection && selection.rangeCount) {
				const range = selection.getRangeAt(0).cloneRange()

				if (target.contains(range.commonAncestorContainer)) {
					return { kind: "range", range }
				}
			}
		}

		return null
	}

	function setTextControlValue(target, value, start, end) {
		const proto =
			target instanceof HTMLTextAreaElement
				? HTMLTextAreaElement.prototype
				: HTMLInputElement.prototype

		const descriptor = Object.getOwnPropertyDescriptor(proto, "value")

		if (descriptor?.set) {
			descriptor.set.call(target, value)
		} else {
			target.value = value
		}

		target.setSelectionRange(start, end)

		target.dispatchEvent(
			new InputEvent("input", {
				bubbles: true,
				inputType: "insertText",
				data: value.slice(start, end),
			}),
		)

		target.dispatchEvent(new Event("change", { bubbles: true }))
	}

	function restoreAndInsert(target, saved, replacement) {
		if (
			saved?.kind === "text-control" &&
			(target instanceof HTMLTextAreaElement ||
				target instanceof HTMLInputElement)
		) {
			const value = target.value
			const start = Math.max(0, Math.min(saved.start, value.length))
			const end = Math.max(start, Math.min(saved.end, value.length))

			const nextValue = value.slice(0, start) + replacement + value.slice(end)

			const caret = start + replacement.length

			setTextControlValue(target, nextValue, caret, caret)
			target.focus()
			return
		}

		if (saved?.kind === "range" && target?.isContentEditable) {
			const selection = window.getSelection()
			if (!selection) return

			selection.removeAllRanges()
			selection.addRange(saved.range)

			const inserted = document.execCommand("insertText", false, replacement)

			if (!inserted) {
				const range = saved.range
				range.deleteContents()

				const textNode = document.createTextNode(replacement)
				range.insertNode(textNode)
				range.setStartAfter(textNode)
				range.collapse(true)

				selection.removeAllRanges()
				selection.addRange(range)

				target.dispatchEvent(
					new InputEvent("input", {
						bubbles: true,
						inputType: "insertText",
						data: replacement,
					}),
				)
			}

			target.focus()
		}
	}

	let busy = false

	async function handlePaste(event) {
		console.info("[Goredacted] Paste event detected:", event.target)

		if (busy) {
			console.warn("[Goredacted] Paste skipped: already processing")
			return
		}

		if (!isEnabledHere()) {
			console.warn("[Goredacted] Paste skipped: domain disabled")
			return
		}

		const target = getEditableTarget(event.target)
		console.info("[Goredacted] Editable target:", target)

		if (!target) {
			console.warn("[Goredacted] Paste skipped: unsupported target")
			return
		}

		const text = event.clipboardData?.getData("text/plain")

		if (text == null || text.length === 0) {
			console.warn("[Goredacted] Paste skipped: no plain text")
			return
		}

		const saved = saveSelection(target)

		if (!saved) {
			console.warn("[Goredacted] Paste skipped: selection unavailable")
			return
		}

		event.preventDefault()
		event.stopPropagation()

		busy = true
		target.setAttribute("data-goredacted-paste-pending", "true")

		try {
			console.info("[Goredacted] Sending text to background")

			const result = await chrome.runtime.sendMessage({
				type: "REDACT_PASTE",
				text,
			})

			console.info("[Goredacted] Background response:", result)

			if (!result?.ok) {
				restoreAndInsert(target, saved, text)

				if (result?.skipped) {
					console.warn("[Goredacted] Background skipped this domain")
				}

				if (result?.error) {
					console.error("[Goredacted] Background error:", result.error)
				}

				return
			}

			restoreAndInsert(target, saved, result.text)
			console.info("[Goredacted] Redacted text inserted")
		} catch (error) {
			console.error("[Goredacted] Paste processing failed:", error)
			restoreAndInsert(target, saved, text)
		} finally {
			target.removeAttribute("data-goredacted-paste-pending")
			busy = false
		}
	}

	document.addEventListener("paste", handlePaste, true)

	chrome.storage.onChanged.addListener((changes, areaName) => {
		if (areaName !== "local") return
		loadSettings().catch((error) => {
			console.error("[Goredacted] Failed to reload settings:", error)
		})
	})

	loadSettings().catch((error) => {
		console.error("[Goredacted] Failed to load settings:", error)
	})
})()
