// ------------------------ General Options ------------------------
const httpEndpointInput = document.getElementById("http-endpoint")
const saveGeneralBtn = document.getElementById("save-general")

// Load saved general options
chrome.storage.sync.get(["httpEndpoint"], (result) => {
	if (result.httpEndpoint) {
		httpEndpointInput.value = result.httpEndpoint
	}
})

// Save general options
saveGeneralBtn.addEventListener("click", () => {
	const endpoint = httpEndpointInput.value.trim()
	chrome.storage.sync.set({ httpEndpoint: endpoint }, () => {
		alert("General options saved!")
	})
})

// ------------------------ Enabled Domains ------------------------
const newEnabledInput = document.getElementById("new-enabled-domain")
const addEnabledBtn = document.getElementById("add-enabled")
const enabledList = document.getElementById("enabled-domain-list")

// ------------------------ Disabled Domains ------------------------
const newDisabledInput = document.getElementById("new-disabled-domain")
const addDisabledBtn = document.getElementById("add-disabled")
const disabledList = document.getElementById("disabled-domain-list")

// ------------------------ Helper Functions ------------------------
function loadDomains() {
	chrome.storage.sync.get(["enabledDomains", "disabledDomains"], (result) => {
		const enabled = result.enabledDomains || []
		const disabled = result.disabledDomains || []

		renderList(enabledList, enabled, (index) => {
			enabled.splice(index, 1)
			chrome.storage.sync.set({ enabledDomains: enabled }, loadDomains)
		})

		renderList(disabledList, disabled, (index) => {
			disabled.splice(index, 1)
			chrome.storage.sync.set({ disabledDomains: disabled }, loadDomains)
		})
	})
}

function renderList(listElement, domains, onDelete) {
	listElement.innerHTML = ""
	domains.forEach((domain, index) => {
		const li = document.createElement("li")
		li.textContent = domain

		const deleteBtn = document.createElement("button")
		deleteBtn.textContent = "Delete"
		deleteBtn.className = "delete-btn"
		deleteBtn.addEventListener("click", () => onDelete(index))

		li.appendChild(deleteBtn)
		listElement.appendChild(li)
	})
}

function addDomain(input, storageKey, listElement) {
	const domain = input.value.trim()
	if (!domain) return

	chrome.storage.sync.get([storageKey], (result) => {
		const domains = result[storageKey] || []
		if (!domains.includes(domain)) {
			domains.push(domain)
			chrome.storage.sync.set({ [storageKey]: domains }, loadDomains)
		}
		input.value = ""
	})
}

// ------------------------ Event Listeners ------------------------
addEnabledBtn.addEventListener("click", () => {
	addDomain(newEnabledInput, "enabledDomains", enabledList)
})

addDisabledBtn.addEventListener("click", () => {
	addDomain(newDisabledInput, "disabledDomains", disabledList)
})

// Initial load
loadDomains()
