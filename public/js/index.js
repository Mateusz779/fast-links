let currentStream = null;
let html5QrcodeScanner = null;
const statusDiv = document.getElementById("status");
const startButton = document.getElementById("startButton");
const stopButton = document.getElementById("stopButton");
const submitUrlButton = document.getElementById("submitUrl");
const manualUrlInput = document.getElementById("manualUrl");

const urlParams = new URLSearchParams(window.location.search);
const idParam = urlParams.get("id");
document.getElementById("idInput").value = idParam || "";

// Funkcja do obsługi i wyświetlania błędów
function handleError(error) {
  console.error("Wystąpił błąd:", error);

  let errorMessage = "Wystąpił nieoczekiwany błąd.";

  if (error instanceof Error) {
    switch (error.name) {
      case "NotAllowedError":
        errorMessage =
          "Dostęp do kamery został zabroniony. Sprawdź ustawienia przeglądarki.";
        break;
      case "NotFoundError":
        errorMessage =
          "Nie znaleziono kamery. Sprawdź czy urządzenie ma kamerę.";
        break;
      case "NotReadableError":
        errorMessage = "Kamera jest obecnie używana przez inną aplikację.";
        break;
      case "OverconstrainedError":
        errorMessage = "Nie znaleziono odpowiedniej kamery. Spróbuj ponownie.";
        break;
      default:
        errorMessage =
          error.message ||
          "Wystąpił problem z kamerą. Spróbuj odświeżyć stronę.";
    }
  }

  statusDiv.textContent = errorMessage;
  return errorMessage;
}

async function sendUrl(url) {
  try {
    statusDiv.textContent = "Wysyłanie danych...";

    const id = document.getElementById("idInput").value.trim();
    const expire = document.getElementById("expireInput").value.trim();
    const body = { value: url };
    if (id) body.id = id;
    if (expire) body.expire = expire;

    const response = await fetch("/api/set", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    });

    if (response.ok) {
      statusDiv.textContent = "Sukces! URL został przesłany.";
      return true;
    } else {
      throw new Error("Błąd podczas wysyłania danych na serwer");
    }
  } catch (error) {
    handleError(error);
    return false;
  }
}

async function stopScanner() {
  try {
    if (html5QrcodeScanner) {
      await html5QrcodeScanner.stop();
      html5QrcodeScanner.clear();
      html5QrcodeScanner = null;
    }

    if (currentStream) {
      const tracks = currentStream.getTracks();
      tracks.forEach((track) => {
        if (track.readyState === "live") {
          track.stop();
        }
      });
      currentStream = null;
    }

    document.querySelector(".zoom-control").style.display = "none";
    startButton.style.display = "block";
    stopButton.style.display = "none";

    const readerElement = document.getElementById("reader");
    if (readerElement) {
      readerElement.innerHTML = "";
    }
  } catch (error) {
    handleError(error);
  }
}

startButton.addEventListener("click", async function () {
  const zoomControl = document.querySelector(".zoom-control");
  const zoomSlider = document.getElementById("zoomSlider");
  const zoomValue = document.querySelector(".zoom-value");

  try {
    // Najpierw zatrzymaj poprzedni skaner
    await stopScanner();

    // Poczekaj chwilę przed ponownym uruchomieniem
    await new Promise((resolve) => setTimeout(resolve, 500));

    // Sprawdź dostępne urządzenia
    const devices = await navigator.mediaDevices.enumerateDevices();
    const videoDevices = devices.filter(
      (device) => device.kind === "videoinput"
    );

    if (videoDevices.length === 0) {
      throw new Error("Nie znaleziono żadnej kamery w urządzeniu");
    }

    // Preferuj tylną kamerę
    const camera =
      videoDevices.find(
        (device) =>
          device.label.toLowerCase().includes("back") ||
          device.label.toLowerCase().includes("tylna")
      ) || videoDevices[videoDevices.length - 1];

    const stream = await navigator.mediaDevices.getUserMedia({
      video: {
        deviceId: camera.deviceId,
        facingMode: { ideal: "environment" },
      },
    });

    currentStream = stream;
    const track = stream.getVideoTracks()[0];

    if (!track) {
      throw new Error("Nie udało się uzyskać dostępu do kamery");
    }

    // Sprawdź możliwości kamery
    const capabilities = track.getCapabilities();
    if (capabilities && capabilities.zoom) {
      zoomControl.style.display = "block";
      zoomSlider.min = capabilities.zoom.min;
      zoomSlider.max = capabilities.zoom.max;
      zoomSlider.step = (capabilities.zoom.max - capabilities.zoom.min) / 100;
      zoomSlider.value = capabilities.zoom.min;
      zoomValue.textContent = `${capabilities.zoom.min}x`;

      zoomSlider.addEventListener("input", async (e) => {
        try {
          const zoomValue = parseFloat(e.target.value);
          await track.applyConstraints({
            advanced: [{ zoom: zoomValue }],
          });
          document.querySelector(
            ".zoom-value"
          ).textContent = `${zoomValue.toFixed(1)}x`;
        } catch (error) {
          console.warn("Nie udało się ustawić zoomu:", error);
        }
      });
    } else {
      zoomControl.style.display = "none";
    }

    html5QrcodeScanner = new Html5Qrcode("reader");
    const qrConfig = {
      fps: 10,
      qrbox: { width: 250, height: 250 },
      aspectRatio: 1.0,
    };

    startButton.style.display = "none";
    stopButton.style.display = "block";
    statusDiv.textContent = "Skaner jest aktywny...";

    await html5QrcodeScanner.start(
      { facingMode: "environment" },
      qrConfig,
      async (decodedText) => {
        await stopScanner();
        if (await sendUrl(decodedText)) {
          manualUrlInput.value = decodedText;
        }
      },
      (error) => {
        // Ignorujemy błędy skanowania, bo pojawiają się często gdy nie ma kodu QR w kadrze
        console.debug("Skanowanie:", error);
      }
    );
  } catch (error) {
    await stopScanner();
    handleError(error);
  }
});

stopButton.addEventListener("click", stopScanner);
window.addEventListener("beforeunload", stopScanner);
document.addEventListener("visibilitychange", () => {
  if (document.hidden && html5QrcodeScanner) {
    stopScanner();
  }
});

submitUrlButton.addEventListener("click", async () => {
  const url = manualUrlInput.value.trim();
  if (!url) {
    statusDiv.textContent = "Proszę wprowadzić URL";
    return;
  }

  try {
    new URL(url); // Sprawdź poprawność URL
    await sendUrl(url);
  } catch (error) {
    if (error instanceof TypeError) {
      statusDiv.textContent = "Wprowadzony URL jest nieprawidłowy";
    } else {
      handleError(error);
    }
  }
});

manualUrlInput.addEventListener("keypress", (e) => {
  if (e.key === "Enter") {
    submitUrlButton.click();
  }
});
