const DEFAULTS = {
  endpoint: "http://127.0.0.1:5678/parse",
  enabledDomains: [],
  disabledDomains: []
};

const endpoint = document.getElementById("endpoint");
const enabled = document.getElementById("enabled");
const disabled = document.getElementById("disabled");
const save = document.getElementById("save");
const status = document.getElementById("status");

function lines(value) {
  return value
    .split(/\r?\n/)
    .map(s => s.trim())
    .filter(Boolean);
}

async function restore() {
  const settings = await chrome.storage.local.get(DEFAULTS);
  endpoint.value = settings.endpoint;
  enabled.value = settings.enabledDomains.join("\n");
  disabled.value = settings.disabledDomains.join("\n");
}

async function saveSettings() {
  let parsedEndpoint;
  try {
    parsedEndpoint = new URL(endpoint.value.trim());
  } catch {
    status.textContent = "Enter a valid endpoint URL.";
    return;
  }

  if (!["http:", "https:"].includes(parsedEndpoint.protocol)) {
    status.textContent = "The endpoint must use HTTP or HTTPS.";
    return;
  }

  await chrome.storage.local.set({
    endpoint: parsedEndpoint.toString(),
    enabledDomains: lines(enabled.value),
    disabledDomains: lines(disabled.value)
  });

  status.textContent = "Saved.";
  setTimeout(() => {
    status.textContent = "";
  }, 1800);
}

save.addEventListener("click", saveSettings);
restore();
