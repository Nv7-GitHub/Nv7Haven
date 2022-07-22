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
export let uid = writable("");

async function load() {
  let val = await fetch(host + "services");
  services.set(await val.json());
}

load();

if (browser) {
  let sse = new EventSource(host + "events?stream=services");
  sse.addEventListener("message", (e) => {
    let dat = JSON.parse(e.data);
    services.set(dat);
  });
}
