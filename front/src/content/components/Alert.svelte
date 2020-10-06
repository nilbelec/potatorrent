<script>
  import { scale } from "svelte/transition";
</script>

<style>
  .alert {
    text-align: center;
    padding: 1rem;
  }
</style>

<div class="alert" in:scale={{}}>
  <slot />
</div>
