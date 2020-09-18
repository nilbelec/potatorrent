<script>
  import Search from "./Search.svelte";
</script>

<style>
  main {
    padding: 5rem 6rem 2rem 6rem;
    transition: padding 0.3s;
  }
  @media only screen and (max-width: 1000px) {
    main {
      padding: 5rem 2rem 2rem 2rem;
    }
  }
</style>

<main>
  <Search />
</main>
