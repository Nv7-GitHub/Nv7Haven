<script lang="ts">
  import { page } from '$app/stores';
  import { host, services, uid } from '$lib/config';
  import { browser } from '$app/env';

  let service = $services.find(service => service.id == $page.params.service)!;
  services.subscribe((dat) => {
    service = dat.find(service => service.id == $page.params.service)!;
  })

  // Connection
  let logs = "";
  if (browser) {
    let socket = new WebSocket("ws" + host.substring(4) + "logs/" + service.id);
    socket.onmessage = (data) => {
      logs += data.data;
    }
  }

  async function start() {
    await fetch(host + "start/" + service.id, {
      method: "POST",
      body: $uid,
    })
  }

  async function stop() {
    await fetch(host + "stop/" + service.id, {
      method: "POST",
      body: $uid,
    })
  }

  async function rebuild() {
    await fetch(host + "rebuild/" + service.id, {
      method: "POST",
      body: $uid,
    })
  }
</script>

<svelte:head>
  <title>Nv7 Status - {service.name}</title>
</svelte:head>

<div class="container">
  <pre class="block mt-3"><code>{logs}</code></pre>

  {#if $uid != ""}
    <div class="text-center">
      <div class="btn-group" role="group">
        <button type="button" class="btn btn-secondary" disabled={service.building} on:click={rebuild}>Rebuild</button>
        <button type="button" class="btn btn-danger" disabled={!service.running} on:click={stop}>Stop</button>
        <button type="button" class="btn btn-success" disabled={service.running} on:click={start}>Start</button>
      </div>
    </div>
  {/if}
</div>

<style>
  .block {
    background-color: #f5f5f5;
    border-radius: 10px;
    padding: 1em;
  }
</style>