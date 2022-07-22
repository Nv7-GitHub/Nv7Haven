<script lang="ts">
  import { page } from '$app/stores';
  import { services } from '$lib/config';

  let service = $services.find(service => service.id == $page.params.service)!;
</script>

<svelte:head>
  <title>Nv7 Status - {service.name}</title>
</svelte:head>

<div class="container">
  <pre class="block mt-3">
    code
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
  }
</style>