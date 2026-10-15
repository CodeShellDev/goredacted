const { DEFAULTS, isDomainEnabled } = GoredactedDomainRules

const statusElement = document.getElementById("status")
const endpointElement = document.getElementById("endpoint")
const healthLabel = document.getElementById("health-label")
const healthDot = document.getElementById("health-dot")
const healthButton = document.getElementById("check-health")

function setPageStatus(message, enabled = null) {
	const dotClass = enabled === true ? "on" : enabled === false ? "off" : ""

	statusElement.replaceChildren()

	const dot = document.createElement("span")
	dot.className = `dot ${dotClass}`

	const label = document.createElement("span")
	label.className = "status-label"
	label.textContent = message

	statusElement.append(dot, label)
}

function setHealthStatus(message, state = "checking") {
	healthLabel.textContent = message
	healthDot.className = `dot ${
		state === "healthy" ? "on" : state === "unhealthy" ? "off" : ""
	}`
}

async function init() {
	try {
		const settings = await chrome.storage.local.get(DEFAULTS)
		endpointElement.textContent = settings.endpoint || ""

		const [tab] = await chrome.tabs.query({
			active: true,
			currentWindow: true,
		})

		if (!tab?.url) {
			setPageStatus("Unable to read active tab URL", false)
		} else {
			let url

			try {
				url = new URL(tab.url)
			} catch {
				setPageStatus("Unsupported page URL", false)
				url = null
			}

			if (url) {
				if (!["http:", "https:"].includes(url.protocol)) {
					setPageStatus("Not available on this page", false)
				} else {
					const enabled = isDomainEnabled(url.hostname, settings)

					setPageStatus(
						enabled ? "Enabled on this site" : "Disabled on this site",
						enabled,
					)
				}
			}
		}
	} catch (error) {
		console.error("[Goredacted] Popup initialization failed:", error)
		setPageStatus("Could not determine page status", false)
	}

	await refreshHealth()
}

async function refreshHealth() {
	healthButton.disabled = true
	setHealthStatus("Checking local service…")

	try {
		const result = await chrome.runtime.sendMessage({
			type: "CHECK_HEALTH",
		})

		if (!result?.ok || !result.healthy) {
			throw new Error(result?.error || "Health check failed")
		}

		setHealthStatus("Local service healthy", "healthy")
		console.log("[Goredacted] Health check passed.")
	} catch (error) {
		console.error("[Goredacted] Health check failed:", error)
		setHealthStatus("Local service unavailable", "unhealthy")
	} finally {
		healthButton.disabled = false
	}
}

healthButton.addEventListener("click", refreshHealth)

document.getElementById("settings").addEventListener("click", () => {
	chrome.runtime.openOptionsPage()
})

init()
