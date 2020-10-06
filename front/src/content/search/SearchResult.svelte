<script>
  import Torrent from "./Torrent.svelte";

  export let torrents = [];
</script>

<style>
  .results-grid {
    display: grid;
    grid-gap: 3rem;
    grid-template-columns: repeat(auto-fit, 172px);
    justify-content: space-evenly;
  }
</style>

<div class="results-grid">
  {#each torrents as torrent}
    <Torrent {torrent} />
  {/each}
</div>
