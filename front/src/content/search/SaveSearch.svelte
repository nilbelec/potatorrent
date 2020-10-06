<script>
  import { fly } from "svelte/transition";
  import Icon from "svelte-awesome";
  import { check, warning, save } from "svelte-awesome/icons";
  import { saveSearch } from "../../api.js";

  export let searchParams = {};
  export let disabled = false;

  let success = false;
  let error = false;

  async function onClick() {
    disabled = true;
    success = false;
    error = false;
    try {
      await saveSearch(searchParams);
      success = true;
      setTimeout(() => {
        success = false;
        disabled = false;
      }, 3000);
    } catch (e) {
      error = true;
      setTimeout(() => {
        error = false;
        disabled = false;
      }, 3000);
    }
  }
</script>

<style>
  button {
    display: block;
    padding: 8px 24px;
    display: flex;
    align-items: center;
  }
  button > span {
    margin-right: 8px;
  }
  .message {
    padding: 8px 16px;
    color: white;
    display: flex;
    align-items: center;
    margin-right: 1rem;
    border-radius: 30px;
  }
  .success {
    background: rgb(8, 182, 17);
  }
  .error {
    background: red;
  }
</style>

{#if error}
  <div class="message error" transition:fly|local={{ duration: 1000, x: 200 }}>
    <span>No se ha podido guardar la búsqueda</span>
    <Icon data={warning} />
  </div>
{:else if success}
  <div
    class="message success"
    transition:fly|local={{ duration: 1000, x: 200 }}>
    <span>Búsqueda guardada</span>
    <Icon data={check} />
  </div>
{/if}
<button {disabled} on:click={onClick}>
  <span>Programar esta búsqueda</span>
  <Icon data={save} />
</button>
