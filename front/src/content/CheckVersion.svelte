<script>
  import Icon from "svelte-awesome";
  import { bell, times } from "svelte-awesome/icons";
  import { scale } from "svelte/transition";
  import { shrink } from "../transitions/shrink.js";
  import { isNewVersion } from "../api.js";

  let close = false;
</script>

<style>
  @keyframes tada {
    0% {
      transform: scale(1);
    }

    10%,
    20% {
      transform: scale(0.9) rotate(-8deg);
    }

    30%,
    50%,
    70% {
      transform: scale(1.3) rotate(8deg);
    }

    40%,
    60% {
      transform: scale(1.3) rotate(-8deg);
    }

    100%,
    80% {
      transform: scale(1) rotate(0);
    }
  }
  .alert {
    padding: 0 1rem;
    height: 60px;
    border: 1px solid rgba(36, 241, 6, 0.46);
    background-color: rgba(7, 149, 66, 0.12156862745098039);
    box-shadow: 0px 0px 2px #259c08;
    color: #0b9908;
    border-color: #c3e6cb;
    display: flex;
    align-items: center;
    font-size: 0.9rem;
  }
  .icon {
    display: flex;
    align-items: center;
  }
  .icon.bell {
    animation: tada 2s linear infinite;
    margin-right: 1rem;
  }
  .icon.close {
    margin-left: auto;
    cursor: pointer;
  }
  a {
    font-weight: bolder;
  }
</style>

{#await isNewVersion() then value}
  {#if value && !close}
    <div class="alert" transition:shrink={{ duration: 600 }}>
      <span class="icon bell">
        <Icon data={bell} scale={1.4} />
      </span>
      <div>
        <strong>¡Yuhuuu!</strong>
        Hay una nueva versión disponible: {value}.
        <a
          href="https://github.com/nilbelec/potatorrent/releases"
          target="_blank">
          Ir a la página de descargas
        </a>
      </div>
      <span class="icon close" on:click={() => (close = true)}>
        <Icon data={times} scale={1.4} />
      </span>
    </div>
  {/if}
{/await}
