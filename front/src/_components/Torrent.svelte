<script>
  import Lazy from "svelte-lazy";
  export let torrent;
</script>

<style>
  .torrent {
    padding: 4px;
    transition: transform 280ms cubic-bezier(0.4, 0, 0.2, 1),
      box-shadow 280ms cubic-bezier(0.4, 0, 0.2, 1);
  }
  .torrent:hover {
    transform: translateY(-4px);
    box-shadow: 0 5px 5px -3px rgba(0, 0, 0, 0.2),
      0 8px 10px 1px rgba(0, 0, 0, 0.14), 0 3px 14px 2px rgba(0, 0, 0, 0.12);
  }
  .quality-cont {
    text-align: center;
    padding-bottom: 4px;
  }
  .quality {
    text-overflow: ellipsis;
    white-space: nowrap;
    font-size: 0.7rem;
    font-weight: bolder;
  }
  .image {
    height: 240px;
    background-repeat: no-repeat;
    background-position: center;
    background-size: contain;
  }
  .name {
    text-align: center;
    padding: 4px;
  }
  .show {
    text-align: center;
    font-size: 0.8rem;
  }
  .show > span {
    display: inline-block;
  }
</style>

<div class="torrent">
  <div class="quality-cont">
    <span class="quality">{torrent.calidad}</span>
  </div>
  <Lazy height={240}>
    <div
      class="image"
      style="background-image:url({'/image?path=' + torrent.imagen})" />
  </Lazy>
  <div class="name">{torrent.name}</div>
  {#if torrent.season}
    <div class="show">
      <span>Temporada {torrent.season}</span>
      <span>
        {#if torrent.singleEpisode}
          Capítulo {torrent.firstEpisode}
        {:else}
          Capítulos del {torrent.firstEpisode} al {torrent.lastEpisode}
        {/if}
      </span>
    </div>
  {/if}
</div>
