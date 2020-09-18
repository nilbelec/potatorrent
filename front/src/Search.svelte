<script>
  import { onMount } from "svelte";
  import SearchForm from "./_components/SearchForm.svelte";
  import SearchResult from "./_components/SearchResult.svelte";
  import SvelteInfiniteScroll from "svelte-infinite-scroll";
  import { BarLoader } from "svelte-loading-spinners";

  import { searchTorrents } from "./api.js";

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
    const searchResult = await searchTorrents(
      searchParams.category && searchParams.category.value,
      searchParams.subCategory && searchParams.subCategory.value,
      searchParams.quality && searchParams.quality.value,
      searchParams.text,
      page
    );
    if (searchResult && searchResult.data && searchResult.data.torrents) {
      newBatch = searchResult.data.torrents;
    } else {
      newBatch = [];
    }
    loading = false;
  }

  onMount(fetchData);

  let loading = false;
  let page = 1;
  let searchParams;
  $: torrents = [...torrents, ...newBatch];
</script>

<style>
  .loading {
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 4rem 0;
  }
</style>

<SearchForm bind:searchParams handleSubmit={search} />
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
  }} />
