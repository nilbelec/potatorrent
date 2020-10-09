<script>
  import { scale } from "svelte/transition";
  import { Route } from "svelte-routing";

  export let path;

  const duration = 280;
  let inProps = { y: -1000, duration: duration };
  let outProps = { duration: duration, y: 1000 };
</script>

<style>
  .animated {
    position: absolute;
    width: 100%;
    padding-bottom: 2rem;
  }
</style>

<Route {path}>
  <div class="animated" in:scale={inProps} out:scale={outProps}>
    <slot />
  </div>
</Route>
