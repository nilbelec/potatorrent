<script>
  import { isNewVersion } from "../../api.js";
  import Icon from "svelte-awesome";
  import { bell, refresh } from "svelte-awesome/icons";
  import Error from "../components/Error.svelte";
</script>

<style>
  .loading {
    padding: 1rem;
    text-align: center;
  }
  .new-version {
    margin-bottom: 1rem;
    background: rgb(58, 166, 58);
    border: 1px solid green;
    color: white;
    padding: 1rem;
    display: flex;
    align-items: center;
  }
  .icon {
    display: flex;
    align-items: center;
    margin-right: 1rem;
  }
  a {
    font-weight: bolder;
  }
</style>

{#await isNewVersion()}
  <div class="loading">
    <Icon data={refresh} spin scale={3} />
  </div>
{:then value}
  <div class="new-version">
    <span class="icon">
      <Icon data={bell} scale={2} />
    </span>
    <div>
      ¡Hay una nueva versión disponible!
      <a
        target="_blank"
        href="https://github.com/nilbelec/potatorrent/releases">
        Descárgala aquí
      </a>
    </div>
  </div>
{:catch error}
  <Error message="No se ha podido comprobar si hay una nueva versión..." />
{/await}
