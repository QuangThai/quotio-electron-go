import { contextBridge } from "electron";
import type { ElectronAPI } from "../types";

const electronAPI: ElectronAPI = {
  platform: process.platform,
  versions: process.versions,
};

contextBridge.exposeInMainWorld("electron", electronAPI);

