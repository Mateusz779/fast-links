let lastData = "";
const urlParams = new URLSearchParams(window.location.search);
const id = urlParams.get("id");
let fetchingData;
let newWindow;
let listening = false;
const startBtn = document.getElementById("startBtn");
const autoListen = document.getElementById("autoListen");
const saveOldLink = document.getElementById("saveOldLink");
const info = document.getElementById("info");
const toast = document.getElementById("toast");

autoListen.addEventListener("click", function () {
  if (autoListen.checked) {
    saveOldLink.checked = true;
    toast.textContent = "Attention! You must allow pop-ups to use this feature!";
    toast.classList.add("show");
    setTimeout(() => {
      toast.classList.remove("show");
    }, 3000);
  }
});

saveOldLink.addEventListener("click", function () {
  if (autoListen.checked && !this.checked) autoListen.checked = false;
});

startBtn.addEventListener("click", function () {
  info.textContent = "Waiting for a new address...";
  listening = !listening;
  if (listening) startBtn.textContent = "Stop listening";
  else {
    startBtn.textContent = "Start listening";
    clearInterval(fetchingData);
    return;
  }

  if (!saveOldLink.checked) {
    lastData = "";
    const a = document.getElementById("url");
    a.href = "";
    a.textContent = "";
  }
  openNewTab();
  fetchingData = setInterval(checkForNewAddress, 1000);
});

function openNewTab() {
  newWindow = window.open("about:blank", "_blank", "");
}

function checkForNewAddress() {
  const apiUrl = id ? `/api/get?id=${encodeURIComponent(id)}` : "/api/get";

  fetch(apiUrl, {
    method: "GET",
    headers: {
      "Cache-Control": "no-cache",
      Pragma: "no-cache",
    },
  })
    .then((response) => {
      if (!response.ok) {
        throw new Error("Brak danych");
      }
      return response.json();
    })
    .then((data) => {
      if (data.data && data.data !== lastData) {
        lastData = data.data;
        console.log("Nowy adres:", data.data);

        // Przejście do nowego adresu z resetowaniem nagłówków
        const url = new URL(data.data);

        const a = document.getElementById("url");
        a.href = url;
        a.textContent = url;

        if (newWindow) {
          console.log("ustawiono nowy adres");
          newWindow.location.href = id
            ? `/api/url?id=${encodeURIComponent(id)}`
            : "/api/url";
        }

        info.textContent =
          "The page has opened; if not, click the link below. Link waiting is disabled.";

        if (!autoListen.checked) {
          clearInterval(fetchingData);
          listening = false;
          document.getElementById("startBtn").textContent = "Start listening";
        } else {
          openNewTab();
        }
      }
    })
    .catch((error) => {
      console.error("Błąd:", error);
    });
}
