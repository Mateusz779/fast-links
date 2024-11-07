let lastData = "";
const urlParams = new URLSearchParams(window.location.search);
const id = urlParams.get("id");
let new_url = "";
let fetchingData;
let newWindow;
document.getElementById("startBtn").addEventListener("click", function () {
  document.getElementById("info").textContent = "Waiting for a new address...";
  newWindow = window.open("about:blank", "_blank", "");
  fetchingData = setInterval(checkForNewAddress, 1000);
});

function checkForNewAddress() {
  // Dodanie parametru id do URL zapytania

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
        console.log(newWindow);
        if (newWindow) {
          console.log("ustawiono nowy adres");
          newWindow.location.href = id
            ? `/api/url?id=${encodeURIComponent(id)}`
            : "/api/url";
        }

        document.getElementById("info").textContent =
          "The page has opened; if not, click the link below. Link waiting is disabled.";
        clearInterval(fetchingData);
      }
    })
    .catch((error) => {
      console.error("Błąd:", error);
    });
}
