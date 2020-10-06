<script>
  import Alert from "./Alert.svelte";

  export let message = "";
</script>

<style>
  .message {
    font-size: 1.5rem;
    color: green;
  }
</style>

<Alert>
  <div class="message">{message}</div>
</Alert>
