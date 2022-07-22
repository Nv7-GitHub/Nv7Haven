<script lang="ts">
  import { host, services, uid } from "$lib/config";

  let password = "";
  let loading = false;
  async function login() {
    loading = true;
    let uid_req = await fetch(host + "auth", {
      method: "POST",
      body: password,
  });
  uid.set(await uid_req.text());
  loading = false;

  let closebtn = document.getElementById("close-btn")!;
  closebtn.click();
}
</script>
<nav class="navbar navbar-dark bg-dark">
  <div class="container-fluid">
    <a class="navbar-brand" href="/">Nv7 Status</a>

    {#if $uid == ""}
      <button class="btn btn-outline-primary" type="button" data-bs-toggle="modal" data-bs-target="#passwordModal">Authenticate</button>
    {/if}
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

<div class="modal fade" id="passwordModal" tabindex="-1" aria-labelledby="exampleModalLabel" aria-hidden="true">
  <div class="modal-dialog">
    <div class="modal-content">
      <div class="modal-header">
        <h5 class="modal-title" id="exampleModalLabel">Authenticate</h5>
        <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
      </div>
      <form>
        <div class="modal-body">
          <label for="password" class="col-form-label">Password: </label>
          <input type="password" id="password" bind:value={password} autocomplete="current-password"/>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" data-bs-dismiss="modal" id="close-btn">Close</button>
          <button type="submit" class="btn btn-primary" on:click|preventDefault={login} disabled={loading}>
            {#if loading}
              <span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>
            {/if}
            Submit
          </button>
        </div>
      </form>
    </div>
  </div>
</div>

<style>
  .center {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
  }
</style>