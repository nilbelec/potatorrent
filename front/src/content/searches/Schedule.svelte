<script>
  import { fade } from "svelte/transition";
  export let schedule;
  import {
    deleteSchedule,
    enableSchedule,
    disableSchedule
  } from "../../api.js";
  import { createEventDispatcher } from "svelte";
  import Icon from "svelte-awesome";
  import { trash, toggleOn, toggleOff } from "svelte-awesome/icons";
  import Toggle from "svelte-toggle";

  const dispatch = createEventDispatcher();

  async function remove() {
    await deleteSchedule(schedule.id);
    dispatch("deleted", schedule);
  }

  async function toggleEnabled({ detail }) {
    if (detail.toggled) await enableSchedule(schedule.id);
    else await disableSchedule(schedule.id);
    dispatch("deleted", schedule);
  }
</script>

<style>
  .schedule {
    height: 100%;
    box-sizing: border-box;
    box-shadow: 0 8px 16px 0 rgba(0, 0, 0, 0.2);
    background-color: white;
    display: flex;
    flex-direction: column;
    padding: 16px;
    transition: box-shadow 400ms, opacity 400ms;
  }
  .schedule.disabled {
    box-shadow: 0 2px 4px 0 rgba(0, 0, 0, 0.2);
    opacity: 0.6;
  }
  .schedule:hover:not(.disabled) {
    box-shadow: 0 8px 24px 0 rgba(0, 0, 0, 0.2);
  }
  .execution {
    margin-top: 0.5rem;
    font-size: 0.8rem;
  }
  .params {
    display: flex;
    flex-direction: column;
    margin-bottom: 1rem;
  }
  .params > span {
    margin: 0.2rem auto;
    padding: 0.3rem 0.5rem;
    background-color: #007bff;
    color: white;
    display: block;
    border-radius: 1rem;
    font-size: 0.75rem;
    font-weight: bold;
    text-align: center;
    box-shadow: 0px 2px 4px black;
  }
  .message {
    margin-top: auto;
    font-size: 0.75rem;
  }
  .last-torrent {
    display: flex;
    padding: 4px;
    height: 148px;
    margin: 0.5rem 0;
    border-radius: 12px;
  }
  .last-torrent img {
    height: 140px;
    border-radius: 8px;
  }
  .last-torrent .info {
    margin-left: 6px;
    margin-right: auto;
    font-size: 0.95rem;
    font-weight: bold;
    display: flex;
    flex-direction: column;
    justify-content: space-evenly;
  }
  .no-last-torrent {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    text-align: center;
    width: 100%;
    font-size: 3rem;
  }
  .options {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  .options .remove {
    background-color: red;
    color: white;
    border: 1px solid red;
    border-radius: 50%;
    width: 32px;
    height: 32px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
  }
  .options .remove:hover {
    background-color: darkred;
    border: 1px solid darkred;
  }
</style>

<div
  class="schedule"
  class:disabled={schedule.disabled}
  in:fade|local={{ duration: 500 }}>
  <div class="params">
    {#if schedule.params.categoriaTexto}
      <span title="Categoría">{schedule.params.categoriaTexto}</span>
    {/if}
    {#if schedule.params.subcategoriaTexto}
      <span title="Subategoría">{schedule.params.subcategoriaTexto}</span>
    {/if}
    {#if schedule.params.calidadTexto}
      <span title="Calidad o tipo">{schedule.params.calidadTexto}</span>
    {/if}
    {#if schedule.params.q}
      <span title="Palabras">"{schedule.params.q}"</span>
    {/if}
    {#if !schedule.params.categoriaTexto && !schedule.params.subcategoriaTexto && !schedule.params.calidadTexto && !schedule.params.q}
      <span>Cualquier torrent</span>
    {/if}
  </div>
  <div class="message">Último torrent encontrado:</div>
  <div class="last-torrent">
    {#if schedule.lastTorrent}
      <div>
        <img
          src={'/image?path=' + schedule.lastTorrent.imagen}
          alt={schedule.lastTorrent.name} />
      </div>
      <div class="info">
        <div title="Nombre del torrent">{schedule.lastTorrent.name}</div>
        {#if schedule.lastTorrent.season}
          <div class="show">
            <span>Temporada {schedule.lastTorrent.season}</span>
            <span>
              {#if schedule.lastTorrent.singleEpisode}
                Capítulo {schedule.lastTorrent.firstEpisode}
              {:else}
                Capítulos del {schedule.lastTorrent.firstEpisode} al {schedule.lastTorrent.lastEpisode}
              {/if}
            </span>
          </div>
        {/if}
      </div>
    {:else}
      <div
        title="Aún no se ha encontrado ningún torrent"
        class="no-last-torrent">
        -
      </div>
    {/if}
  </div>
  <div class="options">
    <Toggle
      title="Habilitar o deshabilitar esta búsqueda"
      toggled={!schedule.disabled}
      toggledColor="#28a745"
      hideLabel
      on:change={toggleEnabled} />
    <button
      title="Eliminar esta búsqueda programada"
      class="remove"
      on:click={remove}>
      <Icon data={trash} scale="1.2" />
    </button>
  </div>
  <div title="Última vez que se ejecutó esta búsqueda" class="execution">
    {schedule.lastExecutionTime || 'No se ha ejecutado nunca'}
  </div>
</div>
