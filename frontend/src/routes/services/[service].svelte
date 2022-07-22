<script lang="ts">
  import { page } from '$app/stores';
  import { host, services } from '$lib/config';
  import { browser } from '$app/env';

  let service = $services.find(service => service.id == $page.params.service)!;

  // Connection
  let logs = "";
  if (browser) {
    let socket = new WebSocket("ws" + host.substring(4) + "logs/" + service.id);
    socket.onmessage = (data) => {
      logs += data.data;
    }
  }
</script>

<svelte:head>
  <title>Nv7 Status - {service.name}</title>
</svelte:head>

<div class="container">
  <pre class="block mt-3">
    {logs}
  </pre>

  <div class="text-center">
    <div class="btn-group" role="group">
      <button type="button" class="btn btn-secondary" disabled={service.building}>Rebuild</button>
      <button type="button" class="btn btn-danger" disabled={!service.running}>Stop</button>
      <button type="button" class="btn btn-success" disabled={service.running}>Start</button>
    </div>
  </div>
</div>

<style>
  .block {
    background-color: #f5f5f5;
    border-radius: 10px;
    padding: 1em 0 1em 0;
  }
</style>