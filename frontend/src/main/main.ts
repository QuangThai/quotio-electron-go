import { ChildProcess, spawn } from "child_process";
import { app, BrowserWindow, Menu, Tray } from "electron";
import fs from "fs";
import http from "http";
import path from "path";

let mainWindow: BrowserWindow | null = null;
let tray: Tray | null = null;
let backendProcess: ChildProcess | null = null;
const isDev = process.env.NODE_ENV === "development" || !app.isPackaged;

function createWindow(): void {
  const preloadPath = path.join(__dirname, "preload.cjs");

  mainWindow = new BrowserWindow({
    width: 1200,
    height: 800,
    webPreferences: {
      preload: preloadPath,
      nodeIntegration: false,
      contextIsolation: true,
      webSecurity: true,
    },
    titleBarStyle: process.platform === "darwin" ? "hiddenInset" : "default",
  });

  if (isDev) {
    mainWindow.loadURL("http://127.0.0.1:3000");
    mainWindow.webContents.openDevTools();
  } else {
    mainWindow.loadFile(path.join(__dirname, "../../dist/index.html"));
  }

  // Ensure window is shown and focused, especially on Windows
  mainWindow.show();
  if (process.platform === "win32") {
    mainWindow.focus();
  }

  mainWindow.on("closed", () => {
    mainWindow = null;
  });
}

function startBackend(): void {
  if (backendProcess) {
    return; // Already running
  }

  // Skip auto-starting backend in development mode
  // Users can manually start it with: cd backend && go run cmd/server/main.go
  if (isDev) {
    console.log("Development mode: Backend not auto-started. Start manually with: cd backend && go run cmd/server/main.go");
    return;
  }

  // In production, use bundled binary
  const platform = process.platform;
  const arch = process.arch === "x64" ? "amd64" : "arm64";
  const ext = platform === "win32" ? ".exe" : "";
  const backendPath = path.join(
    process.resourcesPath,
    "backend",
    `quotio-server-${platform}-${arch}${ext}`
  );

  if (!fs.existsSync(backendPath)) {
    console.error("Backend binary not found:", backendPath);
    return;
  }

  try {
    backendProcess = spawn(backendPath, [], {
      env: { ...process.env },
    });

    if (backendProcess.stdout) {
      backendProcess.stdout.on("data", (data) => {
        console.log(`Backend: ${data}`);
      });
    }

    if (backendProcess.stderr) {
      backendProcess.stderr.on("data", (data) => {
        console.error(`Backend Error: ${data}`);
      });
    }

    backendProcess.on("close", (code) => {
      console.log(`Backend process exited with code ${code}`);
      backendProcess = null;
    });
  } catch (error) {
    console.error("Failed to start backend:", error);
    backendProcess = null;
  }
}

function stopBackend(): void {
  if (backendProcess) {
    backendProcess.kill();
    backendProcess = null;
  }
}

interface ProxyStatusResponse {
  running: boolean;
  port?: number;
}

async function getProxyStatus(): Promise<ProxyStatusResponse> {
  return new Promise((resolve) => {
    const req = http.get('http://localhost:8080/api/proxy/status', (res) => {
      let data = '';
      res.on('data', (chunk) => { data += chunk; });
      res.on('end', () => {
        try {
          resolve(JSON.parse(data));
        } catch {
          resolve({ running: false });
        }
      });
    });
    req.on('error', () => resolve({ running: false }));
    req.setTimeout(2000, () => {
      req.destroy();
      resolve({ running: false });
    });
  });
}

async function toggleProxy(): Promise<Record<string, unknown>> {
  const status = await getProxyStatus();
  const endpoint = status.running ? 'stop' : 'start';
  return new Promise((resolve) => {
    const req = http.request(`http://localhost:8080/api/proxy/${endpoint}`, { method: 'POST' }, (res) => {
      let data = '';
      res.on('data', (chunk) => { data += chunk; });
      res.on('end', () => {
        try {
          resolve(JSON.parse(data));
        } catch {
          resolve({});
        }
      });
    });
    req.on('error', () => resolve({}));
    req.end();
  });
}

function createTray(): void {
  // Determine icon path based on environment
  let iconPath: string;
  if (isDev) {
    // In development, main.cjs is in frontend/dist-electron/
    // So we need to go up to project root then to build/
    iconPath = path.join(__dirname, "../../build/icon.png");
  } else {
    // In production, icon is bundled with resources
    iconPath = path.join(process.resourcesPath, "build", "icon.png");
  }

  try {
    // Check if icon exists
    if (!fs.existsSync(iconPath)) {
      console.warn(`Icon not found at ${iconPath}, tray disabled`);
      // Don't create tray if icon doesn't exist
      return;
    }

    tray = new Tray(iconPath);

    const contextMenu = Menu.buildFromTemplate([
      {
        label: "🔓 Open Dashboard",
        click: () => {
          if (mainWindow) {
            mainWindow.show();
            mainWindow.webContents.send('navigate', '/');
          } else {
            createWindow();
          }
        },
      },
      {
        label: "🚀 Proxy",
        submenu: [
          {
            label: "Start Proxy",
            click: async () => {
              await toggleProxy();
            },
          },
          {
            label: "Stop Proxy",
            click: async () => {
              await toggleProxy();
            },
          },
        ],
      },
      { type: "separator" },
      {
        label: "❌ Quit",
        click: () => {
          stopBackend();
          app.quit();
        },
      },
    ]);

    tray.setToolTip("Quotio - AI Provider Manager");
    tray.setContextMenu(contextMenu);

    tray.on("click", () => {
      if (mainWindow) {
        if (mainWindow.isVisible()) {
          mainWindow.hide();
        } else {
          mainWindow.show();
        }
      } else {
        createWindow();
      }
    });
  } catch (error) {
    console.error("Failed to create tray:", error);
    tray = null;
  }
}

app.whenReady().then(() => {
  startBackend();
  createWindow();
  createTray();

  app.on("activate", () => {
    if (BrowserWindow.getAllWindows().length === 0) {
      createWindow();
    }
  });
});

app.on("window-all-closed", () => {
  if (process.platform !== "darwin") {
    stopBackend();
    app.quit();
  }
});

app.on("before-quit", () => {
  stopBackend();
});

