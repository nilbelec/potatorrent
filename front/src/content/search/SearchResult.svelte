<script>
  import Torrent from "./Torrent.svelte";

  export let torrents = [];
</script>

<div class="results-grid">
  {#each torrents as torrent}
    <Torrent {torrent} />
  {/each}
</div>

<style>
  .results-grid {
    display: grid;
    gap: 3rem;
    grid-template-columns: repeat(auto-fit, 140px);
    justify-content: space-evenly;
  }
  @media only screen and (max-width: 1000px) {
    .results-grid {
      gap: 2rem;
    }
  }
  @media only screen and (max-width: 600px) {
    .results-grid {
      gap: 1.5rem 0.5rem;
    }
  }
</style>
