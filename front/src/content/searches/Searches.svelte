<script>
  import { getSchedules } from "../../api.js";
  import { flip } from "svelte/animate";
  import Info from "../components/Info.svelte";
  import Error from "../components/Error.svelte";
  import { BarLoader } from "svelte-loading-spinners";
  import Schedule from "./Schedule.svelte";
  import { onDestroy } from "svelte";

  let schedules = [];
  let error = false;
  let loading = true;
  let filter = "";
  let timeout;

  const abortTimeout = () => {
    if (timeout) {
      clearTimeout(timeout);
      timeout = undefined;
    }
  };

  const refresh = async () => {
    abortTimeout();
    try {
      schedules = await getSchedules();
    } catch (err) {
      console.error(err);
      error = true;
    }
    loading = false;
    timeout = setTimeout(refresh, 10000);
  };

  const matches = (words, text) => {
    if (!words || !words.length) return true;
    if (!text || !text.trim()) return false;
    const lw = text.trim().toLocaleLowerCase();
    for (let k in words) {
      if (!lw.includes(words[k])) return false;
    }
    return true;
  };

  $: filtered = schedules.filter((s) => {
    if (!filter) return true;
    const words = filter.trim().toLocaleLowerCase().split(" ");
    return (
      matches(words, s.lastTorrent && s.lastTorrent.name) ||
      matches(words, s.params && s.params.categoriaTexto) ||
      matches(words, s.params && s.params.subcategoriaTexto) ||
      matches(words, s.params && s.params.calidadTexto) ||
      matches(words, s.params && s.params.q) ||
      (!s.params.categoriaTexto &&
        !s.params.subcategoriaTexto &&
        !s.params.calidadTexto &&
        !s.params.q &&
        matches(words, "Cualquier torrent"))
    );
  });

  onDestroy(abortTimeout);

  refresh();
</script>

{#if loading}
  <div class="loading">
    <BarLoader />
  </div>
{:else if error}
  <Error message={"Ups! Se ha producido un error al cargar las búsquedas"} />
{:else if schedules.length}
  <input type="search" placeholder="Filtrar" bind:value={filter} />
  {#if filtered.length}
    <div class="schedules">
      {#each filtered as schedule (schedule.id)}
        <div animate:flip={{ duration: 300 }}>
          <Schedule {schedule} on:deleted={refresh} />
        </div>
      {/each}
    </div>
  {:else}
    <Info
      message={"No hay ninguna búsqueda que cumpla ese filtro..."}
      withImage={false}
    />
  {/if}
{:else}
  <Info message={"No tienes ninguna búsqueda programada"} withImage={false} />
{/if}

<style>
  .loading {
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 2rem 0;
  }
  input {
    margin-bottom: 1rem;
  }
  .schedules {
    display: grid;
    grid-gap: 1.4rem;
    grid-template-columns: repeat(auto-fit, 282px);
    justify-content: center;
  }
</style>
