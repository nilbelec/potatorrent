<script>
  import { shrink } from "../../transitions/shrink.js";
  import Error from "../components/Error.svelte";
  import Success from "../components/Success.svelte";
  import { onMount } from "svelte";
  import Icon from "svelte-awesome";
  import { caretDown, caretUp, warning } from "svelte-awesome/icons";
  import { getConfiguration, saveConfiguration } from "../../api.js";

  let disabled = true;
  let errorLoading = false;
  let conf = {};
  const load = async () => {
    try {
      conf = await getConfiguration();
      disabled = false;
    } catch {
      errorLoading = true;
    }
  };
  let errorMessage = "";
  let saved = false;
  async function handleSubmit(event) {
    errorMessage = "";
    saved = false;
    try {
      await saveConfiguration(conf);
      saved = true;
    } catch (e) {
      errorMessage = e;
    }
  }

  onMount(load);

  let showAdvanced = false;
</script>

<style>
  div {
    margin-bottom: 1rem;
  }
  label {
    margin-bottom: 0.4rem;
    display: inline-block;
  }
  small {
    color: rgb(87, 87, 87);
  }
  .show-advanced {
    display: flex;
    align-items: center;
    font-size: 1.2rem;
    justify-content: flex-end;
  }
  .show-advanced > div {
    cursor: pointer;
    display: flex;
    align-items: center;
    font-size: 1.2rem;
    justify-content: flex-end;
    margin-bottom: 0;
  }
  .show-advanced > div:hover {
    color: #40b3ff;
  }
  .show-advanced > div span {
    margin-left: 0.5rem;
  }
  .show-advanced > div .icon {
    transition: transform 300ms;
  }
  .show-advanced > div .icon.rotate {
    transform: rotate(-180deg);
  }
  .advanced {
    border-bottom: 1px solid #feefb3;
    border-radius: 4px;
    box-shadow: 0 3px 1px -2px rgba(0, 0, 0, 0.2);
  }
  .warning {
    padding: 1rem;
    background-color: #feefb3;
    color: #9f6000;
    border-top-left-radius: 4px;
    border-top-right-radius: 4px;
    display: flex;
    align-items: center;
  }
  .warning strong {
    margin-left: 1rem;
  }
  .bottom {
    display: flex;
    align-items: center;
    height: 90px;
  }
  button {
    display: block;
    padding: 8px 24px;
    display: flex;
    align-items: center;
  }
</style>

{#if errorLoading}
  <Error message={'Ups! Se ha producido un error al cargar la configuración'} />
{:else}
  <form on:submit|preventDefault={handleSubmit}>
    <div>
      <label>Carpeta de descarga por defecto</label>
      <input {disabled} bind:value={conf.downloadFolder} required />
      <small>
        En esta carpeta se descargarán los torrents de las búsquedas programadas
      </small>
    </div>
    <div>
      <label>Intervalo de ejecución</label>
      <input
        {disabled}
        bind:value={conf.intervalInMinutes}
        type="number"
        min="1"
        step="1"
        required />
      <small>
        ¿Cada cuántos minutos comprobarán las búsquedas programadas si hay
        torrents nuevos?
      </small>
    </div>
    <div class="show-advanced">
      <div on:click={() => (showAdvanced = !showAdvanced)}>
        <span class="icon" class:rotate={showAdvanced}>
          <Icon data={caretDown} scale={1.5} />
        </span>
        <span>Parámetros avanzados</span>
      </div>
    </div>
    {#if showAdvanced}
      <div class="advanced" transition:shrink|local>
        <div class="warning">
          <Icon data={warning} />
          <strong>¡Ojo!</strong>
          No modifiques estos valores si no sabes lo que estás haciendo
        </div>
        <div>
          <label>Web a examinar</label>
          <input {disabled} bind:value={conf.baseURL} type="url" required />
          <small>URL de la web donde buscar los torrents</small>
        </div>
        <div>
          <label>Puerto de la web</label>
          <input
            {disabled}
            bind:value={conf.port}
            type="number"
            min="8000"
            step="1"
            max="9999"
            required />
          <small>
            Puerto donde se levantará la web. Deberás reiniciar la aplicación
            para que surta efecto.
          </small>
        </div>
      </div>
    {/if}
    <div class="bottom">
      <button {disabled} type="submit">Guardar Cambios</button>
      {#if errorMessage}
        <Error noImage={true} message={errorMessage} />
      {/if}
      {#if saved}
        <Success message={'Cambios Guardados'} />
      {/if}
    </div>
  </form>
{/if}
