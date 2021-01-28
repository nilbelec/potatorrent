<script>
  import Alert from "./Alert.svelte";

  export let message = "Ups! Se ha producido un error";
  export let noImage = false;
</script>

<Alert>
  <div class="message">{message}</div>
  {#if !noImage}
    <img src="/public/img/error.png" alt="No results found" />
  {/if}
</Alert>

<style>
  .message {
    font-size: 1.5rem;
    color: red;
  }
  img {
    margin: auto;
    margin-top: 1rem;
  }
</style>
