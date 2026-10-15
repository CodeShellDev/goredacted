;(() => {
	const DEFAULTS = {
		endpoint: "http://127.0.0.1:5678/parse",
		healthEndpoint: "http://127.0.0.1:5678/health",
		enabledDomains: [],
		disabledDomains: [],
	}

	function normalizeHostname(hostname) {
		return String(hostname).toLowerCase().replace(/\.$/, "")
	}

	function domainMatches(hostname, pattern) {
		hostname = normalizeHostname(hostname)
		pattern = String(pattern).trim().toLowerCase().replace(/\.$/, "")

		if (!pattern) return false
		if (pattern === "*") return true

		if (pattern.startsWith("*.")) {
			return hostname.endsWith("." + pattern.slice(2))
		}

		if (pattern.startsWith("*")) {
			return hostname.endsWith(pattern.slice(1))
		}

		return hostname === pattern
	}

	function isDomainEnabled(hostname, settings) {
		const disabled = settings.disabledDomains || []
		const enabled = settings.enabledDomains || []

		if (disabled.some((p) => domainMatches(hostname, p))) return false
		if (enabled.length === 0) return true

		return enabled.some((p) => domainMatches(hostname, p))
	}

	globalThis.GoredactedDomainRules = Object.freeze({
		DEFAULTS,
		normalizeHostname,
		domainMatches,
		isDomainEnabled,
	})
})()
