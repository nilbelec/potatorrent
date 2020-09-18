<script>
  import Torrent from "./Torrent.svelte";

  export let torrents = [];
</script>

<style>
  div {
    display: grid;
    grid-gap: 2rem;
    grid-template-columns: repeat(auto-fit, minmax(170px, 1fr));
  }
</style>

<div>
  {#each torrents as torrent}
    <Torrent {torrent} />
  {/each}
</div>
