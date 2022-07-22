<script>
import { services } from "$lib/config";
</script>
<nav class="navbar navbar-dark bg-dark">
  <div class="container-fluid">
    <a class="navbar-brand" href="/">Nv7 Status</a>
  </div>
</nav>

{#if $services.length > 0}
  <slot></slot>
{:else}
  <div class="center">
    <div class="spinner-border" role="status">
      <span class="visually-hidden">Loading...</span>
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