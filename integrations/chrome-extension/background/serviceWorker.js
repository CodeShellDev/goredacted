importScripts("../shared/domainRules.js")

const { DEFAULTS, isDomainEnabled } = GoredactedDomainRules

async function getSettings() {
	return chrome.storage.local.get(DEFAULTS)
}

async function checkHealth(settings) {
	const healthEndpoint = String(settings.healthEndpoint || "").trim()

	if (!healthEndpoint) {
		throw new Error("No goredacted health endpoint is configured.")
	}

	const response = await fetch(healthEndpoint, {
		method: "GET",
		signal: AbortSignal.timeout(3000),
	})

	if (!response.ok) {
		throw new Error(`goredacted health check failed: HTTP ${response.status}`)
	}

	const text = (await response.text()).trim()

	if (text !== "ok") {
		throw new Error("goredacted health check returned an unexpected response.")
	}

	return true
}

chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
	if (!message || !["REDACT_PASTE", "CHECK_HEALTH"].includes(message.type)) {
		return
	}

	;(async () => {
		try {
			const settings = await getSettings()

			if (message.type === "CHECK_HEALTH") {
				await checkHealth(settings)
				return { ok: true, healthy: true }
			}

			const senderUrl = sender.tab?.url || ""
			const url = new URL(senderUrl)

			if (!["http:", "https:"].includes(url.protocol)) {
				throw new Error("This page is not an HTTP(S) page.")
			}

			if (!isDomainEnabled(url.hostname, settings)) {
				return { ok: false, skipped: true }
			}

			// Check the local service before sending clipboard content.
			await checkHealth(settings)

			const endpoint = String(settings.endpoint || "").trim()
			const endpointUrl = new URL(endpoint)

			if (!["http:", "https:"].includes(endpointUrl.protocol)) {
				throw new Error("The goredacted endpoint must use HTTP or HTTPS.")
			}

			const response = await fetch(endpointUrl.toString(), {
				method: "POST",
				headers: { "Content-Type": "text/plain; charset=utf-8" },
				body: String(message.text ?? ""),
				signal: AbortSignal.timeout(30000),
			})

			const responseText = await response.text()

			if (!response.ok) {
				throw new Error(
					`goredacted returned HTTP ${response.status}: ${responseText.slice(0, 300)}`,
				)
			}

			return { ok: true, text: responseText }
		} catch (error) {
			console.error("[Goredacted] Request failed:", error)

			return {
				ok: false,
				error: error instanceof Error ? error.message : String(error),
			}
		}
	})().then(sendResponse)

	return true
})
