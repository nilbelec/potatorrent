<script>
  import Select from "svelte-select";
  import Icon from "svelte-awesome";
  import { search } from "svelte-awesome/icons";
  import { searchOptions, searchSubCategories } from "../../api.js";
  import { searchParamsStore, searchOptsStore } from "../../stores.js";
  import { onMount } from "svelte";
  import { get } from "svelte/store";

  export let disabled = false;

  let params = get(searchParamsStore);
  let selectedCategory = params.category && {
    value: params.category,
    label: params.categoryLabel,
  };
  let selectedSubCategory = params.subCategory && {
    value: params.subCategory,
    label: params.subCategoryLabel,
  };
  let selectedQuality = params.quality && {
    value: params.quality,
    label: params.qualityLabel,
  };
  let text = params.text;

  let categories = [];
  let subCategories = [];
  let qualities = [];

  const load = async () => {
    try {
      let opts = get(searchOptsStore);

      if (!opts || !opts.categories) {
        opts = await searchOptions();
        searchOptsStore.set(opts);
      }

      categories = opts.categories;
      qualities = opts.qualities;
      if (params.category) {
        await prepareSubCategories(opts, params.category);
      }
    } catch {}
  };

  async function prepareSubCategories(opts, cat) {
    if (!opts.subCategories) opts.subCategories = {};
    if (!opts.subCategories[cat]) {
      opts.subCategories[cat] = await searchSubCategories(cat);
      searchOptsStore.set(opts);
    }
    subCategories = opts.subCategories[cat];
  }

  async function onSelectCategory(e) {
    selectedCategory = e.detail;
    selectedSubCategory = undefined;
    updateStore();

    let opts = get(searchOptsStore);
    await prepareSubCategories(opts, selectedCategory.value);
  }

  function updateStore() {
    searchParamsStore.set({
      category: selectedCategory && selectedCategory.value,
      categoryLabel: selectedCategory && selectedCategory.label,
      subCategory: selectedSubCategory && selectedSubCategory.value,
      subCategoryLabel: selectedSubCategory && selectedSubCategory.label,
      quality: selectedQuality && selectedQuality.value,
      qualityLabel: selectedQuality && selectedQuality.label,
      text,
    });
  }

  onMount(load);
</script>

<form on:submit|preventDefault>
  <div>
    <Select
      items={categories}
      selectedValue={selectedCategory}
      on:select={onSelectCategory}
      on:clear={() => {
        selectedCategory = undefined;
        selectedSubCategory = undefined;
        updateStore();
      }}
      isDisabled={disabled || categories.length == 0}
      placeholder="Filtrar por categoría"
    />
  </div>
  <div>
    <Select
      items={subCategories}
      selectedValue={selectedSubCategory}
      on:select={(e) => {
        selectedSubCategory = e.detail;
        updateStore();
      }}
      on:clear={() => {
        selectedSubCategory = undefined;
        updateStore();
      }}
      isDisabled={disabled ||
        selectedCategory == undefined ||
        subCategories.length == 0}
      placeholder="Filtrar por subcategoría"
    />
  </div>
  <div>
    <Select
      items={qualities}
      selectedValue={selectedQuality}
      on:select={(e) => {
        selectedQuality = e.detail;
        updateStore();
      }}
      on:clear={() => {
        selectedQuality = undefined;
        updateStore();
      }}
      isDisabled={disabled || qualities.length == 0}
      placeholder="Filtrar por calidad o tipo"
    />
  </div>
  <div class="words">
    <input
      {disabled}
      bind:value={text}
      on:search={updateStore}
      type="search"
      placeholder="Filtrar por palabras"
    />
    <button type="button" on:click={updateStore}>
      <Icon data={search} />
    </button>
  </div>
</form>

<style>
  form {
    display: grid;
    grid-template-columns: 1fr 1fr 1fr;
    column-gap: 10px;
    row-gap: 10px;
    margin-bottom: 1rem;
  }
  div.words {
    grid-column-start: 1;
    grid-column-end: 4;
    display: flex;
    border: 1px solid #d8dbdf;
    border-radius: 3px;
  }
  div.words input {
    border: 0;
    padding-right: 0;
  }
  div.words button {
    border: 0;
    box-shadow: none;
    width: 48px;
    color: inherit;
  }
  @media only screen and (max-width: 500px) {
    div {
      grid-column-start: 1;
      grid-column-end: 4;
    }
  }
</style>
