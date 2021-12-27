<script>
  import AnimatedRoute from "./components/AnimatedRoute.svelte";
  import Searches from "./searches/Searches.svelte";
  import Configuration from "./configuration/Configuration.svelte";
  import Search from "./search/Search.svelte";
  import CheckVersion from "./CheckVersion.svelte";
</script>

<main>
  <CheckVersion />
  <div class="wrapper">
    <AnimatedRoute path="/">
      <Search />
    </AnimatedRoute>
    <AnimatedRoute path="/searches">
      <Searches />
    </AnimatedRoute>
    <AnimatedRoute path="/configuration">
      <Configuration />
    </AnimatedRoute>
  </div>
</main>

<style>
  main {
    padding: 5rem 6rem 4rem 6rem;
    transition: padding 0.3s;
    margin-bottom: 4rem;
  }
  .wrapper {
    position: relative;
  }
  @media only screen and (max-width: 1000px) {
    main {
      padding-left: 2rem;
      padding-right: 2rem;
    }
  }
  @media (max-width: 800px) {
    main {
      padding-left: 0.7rem;
      padding-right: 0.7rem;
    }
  }
  @media only screen and (max-width: 700px) {
    main {
      padding-top: 4rem;
    }
  }
</style>
