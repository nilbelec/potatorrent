<script>
  import { onDestroy } from "svelte";
  import SaveSearch from "./SaveSearch.svelte";
  import SearchForm from "./SearchForm.svelte";
  import Info from "../components/Info.svelte";
  import Error from "../components/Error.svelte";
  import SearchResult from "./SearchResult.svelte";
  import SvelteInfiniteScroll from "svelte-infinite-scroll";
  import { BarLoader } from "svelte-loading-spinners";
  import { searchParamsStore } from "../../stores.js";

  import { searchTorrents } from "../../api.js";

  let torrents = [];
  let newBatch = [];

  async function search() {
    page = 1;
    torrents = [];
    newBatch = [];
    await fetchData();
  }

  async function fetchData() {
    loading = true;
    error = false;
    try {
      const searchResult = await searchTorrents(
        searchParams.category,
        searchParams.subCategory,
        searchParams.quality,
        searchParams.text,
        page
      );
      if (searchResult && searchResult.data && searchResult.data.torrents) {
        newBatch = searchResult.data.torrents;
      } else {
        newBatch = [];
      }
    } catch (err) {
      error = true;
    }
    loading = false;
  }

  let loading = false;
  let error = false;
  let page = 1;
  let searchParams;
  $: torrents = [...torrents, ...newBatch];

  const unsubscribe = searchParamsStore.subscribe((params) => {
    searchParams = params;
    search();
  });

  onDestroy(unsubscribe);
</script>

<SearchForm />
<div class="save-search">
  <SaveSearch {searchParams} disabled={loading} />
</div>

{#if error}
  <Error message={"Ups! Se ha producido un error al realizar la búsqueda..."} />
{:else if !loading && !torrents.length}
  <Info message={"No se han encontrado resultados"} />
{:else}
  <SearchResult {torrents} />
  {#if loading}
    <div class="loading">
      <BarLoader />
    </div>
  {/if}

  <SvelteInfiniteScroll
    hasMore={newBatch.length}
    window={true}
    threshold={300}
    on:loadMore={() => {
      page++;
      fetchData();
    }}
  />
{/if}

<style>
  .loading {
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 2rem 0;
  }
  .save-search {
    margin-bottom: 1rem;
    display: flex;
    justify-content: flex-end;
  }
</style>
