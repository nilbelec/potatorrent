<script>
  import { fade, blur } from "svelte/transition";
  import { Route } from "svelte-routing";

  export let path;

  const duration = 280;
  let inProps = { y: -1000, delay: duration };
  let outProps = { duration: duration, y: 1000 };
</script>

<Route {path}>
  <div class="animated" in:fade={inProps} out:fade={outProps}>
    <slot />
  </div>
</Route>

<style>
  .animated {
    position: absolute;
    width: 100%;
    padding-bottom: 2rem;
  }
</style>
