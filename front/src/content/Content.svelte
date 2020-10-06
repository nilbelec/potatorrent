<script>
  import AnimatedRoute from "./components/AnimatedRoute.svelte";
  import Searches from "./searches/Searches.svelte";
  import Configuration from "./configuration/Configuration.svelte";
  import Search from "./search/Search.svelte";
  import NewVersion from "./update/NewVersion.svelte";

  let wrapper;
  const duration = 280;
  let inProps = { y: -1000, duration: duration };
  let outProps = { duration: duration, y: 1000 };
</script>

<style>
  main {
    padding: 6rem 6rem 0 6rem;
    transition: padding 0.3s;
    width: 100%;
    height: 100%;
    box-sizing: border-box;
  }
  .wrapper {
    position: relative;
    width: 100%;
    height: 100%;
  }
  @media only screen and (max-width: 1000px) {
    main {
      padding-left: 2rem;
      padding-right: 2rem;
    }
  }
</style>

<main>
  <div class="wrapper" bind:this={wrapper}>
    <AnimatedRoute path="/">
      <Search />
    </AnimatedRoute>
    <AnimatedRoute path="/searches">
      <Searches />
    </AnimatedRoute>
    <AnimatedRoute path="/configuration">
      <Configuration />
    </AnimatedRoute>
    <AnimatedRoute path="/update">
      <NewVersion />
    </AnimatedRoute>
  </div>
</main>
