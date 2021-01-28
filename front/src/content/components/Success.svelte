<script>
  import Alert from "./Alert.svelte";

  export let message = "";
</script>

<Alert>
  <div class="message">{message}</div>
</Alert>

<style>
  .message {
    font-size: 1.5rem;
    color: green;
  }
</style>
