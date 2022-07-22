import { writable, type Writable } from "svelte/store";
import { browser } from '$app/env';

export const host = "http://localhost:8000/";

export type Service = {
  id: string,
  name: string,
  running: boolean,
  building: boolean,
}

export let services: Writable<Service[]> = writable([]);

if (browser) {
  let sse = new EventSource(host + "events?stream=services");
  sse.addEventListener("message", (e) => {
    services.set(JSON.parse(e.data));
  });
}

async function load() {
  let val = await fetch(host + "services");
  services.set(await val.json());
}

load();