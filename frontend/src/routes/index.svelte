<script>
import { services } from "$lib/config";
</script>

<svelte:head>
  <title>Nv7 Status</title>
</svelte:head>

<div class="container mt-2">
  <div class="list-group">
    {#each $services as service}
      <a type="button" href={"/services/" + service.id} class="list-group-item list-group-item-action">{service.name}</a>  
    {/each}
  </div>
</div>
