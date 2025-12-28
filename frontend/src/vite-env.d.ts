/// <reference types="vite/client" />

interface Window {
  electron: {
    platform: string;
    versions: {
      node?: string;
      chrome?: string;
      electron?: string;
    };
  };
}

