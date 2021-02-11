<script>
  import DownloadButton from "./DownloadButton.svelte";
  import Lazy from "svelte-lazy";
  import Icon from "svelte-awesome";
  import { info } from "svelte-awesome/icons";

  export let torrent;

  let keepOptions = false;
</script>

<div class="torrent">
  <div class="quality">{torrent.calidad}</div>
  <Lazy height={250}>
    <div
      class="image"
      style="background-image:url({'/image?path=' + torrent.imagen})"
    >
      <div class="download-opts" class:keep={keepOptions}>
        <div class="opts-top">
          <div class="date" title="Fecha de publicación">{torrent.date}</div>
          <div class="info" title={torrent.fullName}>
            <Icon data={info} />
          </div>
        </div>
        <DownloadButton bind:show={keepOptions} {torrent} />
        <div class="size" title="Tamaño del contenido">
          {torrent.torrentSize}
        </div>
      </div>
    </div>
  </Lazy>
  <div class="bottom">
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
</div>

<style>
  .torrent {
    display: flex;
    flex-direction: column;
    border-radius: 2px;
  }

  .quality {
    text-align: center;
    padding-bottom: 4px;
    text-overflow: ellipsis;
    white-space: nowrap;
    font-size: 0.7rem;
    font-weight: bolder;
  }
  .image {
    height: 250px;
    background-repeat: no-repeat;
    background-position: center;
    background-size: cover;
    position: relative;
    border-radius: 9px;
  }
  .download-opts {
    position: absolute;
    top: 0;
    left: 0;
    height: 100%;
    width: 100%;
    opacity: 0;
    background-color: rgba(0, 0, 0, 0.7);
    transition: opacity cubic-bezier(0.48, 0.21, 1, 0.79) 280ms;
    display: flex;
    justify-content: center;
    align-items: center;
    flex-direction: column;
    border-radius: 9px;
  }
  .image:hover .download-opts,
  .image .download-opts.keep {
    opacity: 1;
  }
  .download-opts .size {
    color: white;
    margin-bottom: 20px;
    background-color: #28a745;
    padding: 4px 8px;
    border-radius: 4px;
    font-size: 0.8rem;
    font-weight: bolder;
  }
  .download-opts .opts-top {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    color: white;
    align-items: center;
    width: 100%;
    padding: 0.5rem 0.5rem 0 0.5rem;
  }
  .download-opts .opts-top .date {
    font-size: 0.8rem;
    font-weight: bolder;
  }
  .download-opts .opts-top .info {
    background-color: rgb(0, 192, 255);
    border-radius: 12px;
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 1rem;
    cursor: help;
  }
  .bottom {
    display: flex;
    flex-direction: column;
    height: 100%;
    justify-content: space-between;
  }
  .name {
    color: #ff3e00;
    text-align: center;
    padding: 4px 0px 6px 0px;
  }
  .show {
    text-align: center;
    font-size: 0.8rem;
  }
  .show > span {
    display: block;
    text-align: center;
  }
</style>
