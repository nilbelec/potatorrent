<script>
  import Alert from "./Alert.svelte";

  export let message = "";
  export let withImage = true;
</script>

<Alert>
  <div class="message">{message}</div>
  {#if withImage}
    <img src="/public/img/no_results.png" alt="No results found" />
  {/if}
</Alert>

<style>
  .message {
    font-size: 1.5rem;
  }
  img {
    margin: auto;
    margin-top: 1rem;
  }
</style>
