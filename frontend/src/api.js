import config from "@/config";
import * as utils from "@/utils";

export async function getServices(params) {
  if (config.desktop_mode) {
    const m = await import("@/../wailsjs/go/api/APIC");
    return m["Services"](params);
  } else {
    let data = await utils.$get(`/api/services/?${params}`);
    return data;
  }
}
export async function createService(params) {
  if (config.desktop_mode) {
    const m = await import("@/../wailsjs/go/api/APIC");
    let encodedSettings = JSON.stringify(params.moduleSettings);
    params.moduleSettings = encodedSettings;
    return m["CreateService"](params);
  } else {
    let data = await utils.$post(`/api/services/create/`, params);
    return data;
  }
}

export async function getModules() {
  if (config.desktop_mode) {
    const m = await import("@/../wailsjs/go/api/APIC");
    return m["Modules"]();
  } else {
    let data = await utils.$get(`/api/modules/`);
    return data;
  }
}

export async function getNetworks() {
  if (config.desktop_mode) {
    const m = await import("@/../wailsjs/go/api/APIC");
    return m["Networks"]();
  } else {
    let data = await utils.$get(`/api/networks/`);
    return data;
  }
}

export async function getSessions(params) {
  if (config.desktop_mode) {
    const m = await import("@/../wailsjs/go/api/APIC");
    return m["Sessions"](params);
  } else {
    let data = await utils.$get(`/api/sessions/?${params}`);
    return data;
  }
}

export async function getSessionEvents(params) {
  if (config.desktop_mode) {
    const m = await import("@/../wailsjs/go/api/APIC");
    return m["SessionEvents"](params);
  } else {
    let data = await utils.$get(`/api/session/view/?${params}`);
    return data;
  }
}

export async function getTunnels(params) {
  if (config.desktop_mode) {
    const m = await import("@/../wailsjs/go/api/APIC");
    return m["Tunnels"](params);
  } else {
    let data = await utils.$get(`/api/tunnels/?${params}`);
    return data;
  }
}

export async function getFiles(params) {
  if (config.desktop_mode) {
    const m = await import("@/../wailsjs/go/api/APIC");
    return m["Files"](params);
  } else {
    let data = await utils.$get(`/api/files/?${params}`);
    return data;
  }
}
export async function getFileTypes() {
  if (config.desktop_mode) {
    const m = await import("@/../wailsjs/go/api/APIC");
    return m["FileTypes"]();
  } else {
    let data = await utils.$get(`/api/files/types/`);
    return data;
  }
}

export async function getAbout() {
  if (config.desktop_mode) {
    const m = await import("@/../wailsjs/go/api/APIC");
    return m["About"]();
  } else {
    let data = await utils.$get(`/api/about/`);
    return data;
  }
}
