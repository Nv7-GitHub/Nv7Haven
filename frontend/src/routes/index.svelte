<script>
import { services } from "$lib/config";
import Service from "./services/[service].svelte";
</script>

<svelte:head>
  <title>Nv7 Status</title>
</svelte:head>

{#if $services.length == 0}
  <div class="center">
    <div class="spinner-border" role="status">
      <span class="visually-hidden">Loading...</span>
    </div>
  </div>
{:else}
  <div class="container mt-2">
    <div class="list-group">
      {#each $services as service}
        <a type="button" href={"/services/" + service.id} class="list-group-item list-group-item-action">{service.name}</a>  
      {/each}
    </div>
  </div>
{/if}

<style>
  .center {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
  }
</style>