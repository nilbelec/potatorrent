<script>
  import Select from "svelte-select";
  import { searchOptions, searchSubCategories } from "../../api.js";
  import { searchParamsStore, searchOptsStore } from "../../stores.js";
  import { onMount } from "svelte";
  import { get } from "svelte/store";

  export let disabled = false;

  let category = undefined;
  let subCategory = undefined;
  let quality = undefined;
  let text = "";

  let categories = [];
  let subCategories = [];
  let qualities = [];

  const load = async () => {
    try {
      let opts = get(searchOptsStore);
      let params = get(searchParamsStore);

      if (!opts || !opts.categories) {
        opts = await searchOptions();
        searchOptsStore.set(opts);
      }

      categories = opts.categories;
      category = params.category && {
        value: params.category,
        label: params.categoryLabel,
      };

      qualities = opts.qualities;
      quality = params.quality && {
        value: params.quality,
        label: params.qualityLabel,
      };

      if (params.category) {
        await prepareSubCategories(opts, params.category);
      }
      subCategory = params.subCategory && {
        value: params.subCategory,
        label: params.subCategoryLabel,
      };
      text = params.text;
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
    subCategory = undefined;
    let opts = get(searchOptsStore);
    await prepareSubCategories(opts, e.detail.value);
  }

  function submit() {
    searchParamsStore.set({
      category: category && category.value,
      categoryLabel: category && category.label,
      subCategory: subCategory && subCategory.value,
      subCategoryLabel: subCategory && subCategory.label,
      quality: quality && quality.value,
      qualityLabel: quality && quality.label,
      text,
    });
  }

  onMount(load);
</script>

<form on:submit|preventDefault={submit}>
  <Select
    items={categories}
    bind:selectedValue={category}
    on:select={onSelectCategory}
    on:clear={() => (subCategory = undefined)}
    isDisabled={disabled || categories.length == 0}
    placeholder="Filtrar por categoría"
  />
  <Select
    items={subCategories}
    bind:selectedValue={subCategory}
    isDisabled={disabled || category == undefined || subCategories.length == 0}
    placeholder="Filtrar por subcategoría"
  />
  <Select
    items={qualities}
    bind:selectedValue={quality}
    isDisabled={disabled || qualities.length == 0}
    placeholder="Filtrar por calidad o tipo"
  />
  <input
    {disabled}
    value={text}
    type="search"
    placeholder="Filtrar por palabras"
  />
  <button {disabled} type="submit">Buscar</button>
</form>

<style>
  form {
    display: grid;
    grid-template-columns: 1fr 1fr 1fr;
    column-gap: 10px;
    row-gap: 10px;
    margin-bottom: 1rem;
  }
  input {
    grid-column-start: 1;
    grid-column-end: 4;
  }
  button {
    grid-column-start: 1;
    grid-column-end: 4;
    height: 42px;
  }
</style>
