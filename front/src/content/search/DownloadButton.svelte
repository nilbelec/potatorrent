<script>
  import {
    getDownloadInfo,
    getDownloadLink,
    downloadOnFolder
  } from "../../api.js";
  import { copyToClipBoard } from "../../utils.js";
  import Icon from "svelte-awesome";
  import {
    refresh,
    download,
    check,
    folderOpen,
    copy,
    warning,
    key
  } from "svelte-awesome/icons";

  export let torrent;
  export let show = false;

  let downloadInfo;
  let doingRequest = false;
  let error = false;
  let downloaded = false;
  let copied = false;
  let noPassword = false;

  $: show = doingRequest || error || downloaded || copied || noPassword;

  async function downloadTorrent() {
    doingRequest = true;
    downloaded = false;
    if (!downloadInfo) {
      await tryDownloadInfo();
    }
    if (!downloadInfo || !downloadInfo.url) {
      error = true;
    } else {
      const url = getDownloadLink(downloadInfo.url);
      window.location.assign(url);
      downloaded = true;
    }
    doingRequest = false;
  }

  async function downloadTorrentOnFolder() {
    doingRequest = true;
    downloaded = false;
    await tryDownloadOnFolder();
    if (!downloadInfo || !downloadInfo.url) error = true;
    else downloaded = true;
    doingRequest = false;
  }

  async function copyDownloadUrl() {
    doingRequest = true;
    copied = false;
    if (!downloadInfo) {
      await tryDownloadInfo();
    }
    if (!downloadInfo || !downloadInfo.url) {
      error = true;
      doingRequest = false;
    } else {
      await copyToClipBoard(downloadInfo.url);
      copied = true;
      doingRequest = false;
    }
  }
  async function copyPassword() {
    doingRequest = true;
    copied = false;
    noPassword = false;
    if (!downloadInfo) {
      await tryDownloadInfo();
    }
    if (!downloadInfo || !downloadInfo.url) {
      error = true;
    } else if (!downloadInfo.password) {
      noPassword = true;
    } else {
      await copyToClipBoard(downloadInfo.password);
      copied = true;
    }
    doingRequest = false;
  }

  async function tryDownloadInfo() {
    error = false;
    try {
      downloadInfo = await getDownloadInfo(torrent);
    } catch {
      error = true;
    }
  }

  async function tryDownloadOnFolder() {
    error = false;
    try {
      downloadInfo = await downloadOnFolder(torrent);
    } catch {
      error = true;
    }
  }
</script>

<style>
  .container {
    margin: auto;
    display: flex;
    flex-direction: column;
    align-items: center;
    color: white;
  }
  button {
    color: white;
    border: 0;
    font-weight: bolder;
    margin: 0;
    cursor: pointer;
    background-color: red;
    font-size: 1rem;
    padding: 0.6rem 1rem;
  }
  button:hover {
    background-color: darkred;
  }
  button:disabled {
    background-color: rgba(255, 0, 0, 0.4);
    cursor: not-allowed;
  }

  .download-btn {
    border-radius: 30px;
  }

  .buttons {
    display: flex;
    align-items: center;
    margin-top: 2rem;
  }

  .buttons > button {
    font-size: 0.9rem;
    background-color: rgb(119, 119, 119);
    padding: 0.6rem;
  }
  .buttons > button:first-child {
    padding-left: 1rem;
    border-top-left-radius: 20px;
    border-bottom-left-radius: 20px;
  }
  .buttons > button:last-child {
    padding-right: 1rem;
    border-top-right-radius: 20px;
    border-bottom-right-radius: 20px;
  }
  .buttons > button:hover {
    background-color: darkgray;
  }
  .buttons > button:disabled {
    background-color: rgba(119, 119, 119, 0.4);
    cursor: not-allowed;
  }

  .error,
  .success {
    margin: auto;
    display: flex;
    flex-direction: column;
    align-items: center;
  }
  .error {
    color: orange;
  }
  .success {
    color: rgb(0, 217, 0);
  }
  .error > div,
  .success > div {
    margin-top: 0.5rem;
    margin-bottom: 1rem;
    text-align: center;
  }
</style>

{#if doingRequest}
  <div class="container">
    <Icon data={refresh} spin scale={2} />
  </div>
{:else if error}
  <div class="error">
    <Icon data={warning} scale={3} />
    <div>No encontrado...</div>
    <button
      on:click={() => {
        downloadInfo = undefined;
        error = false;
      }}>
      OK
    </button>
  </div>
{:else if downloaded || copied || noPassword}
  <div class="success">
    <Icon data={check} scale={3} />
    {#if downloaded}
      <div>Descargado</div>
    {:else if copied}
      <div>Copiado</div>
    {:else if noPassword}
      <div>No necesita contraseña</div>
    {/if}
    <button
      on:click={() => {
        downloaded = copied = noPassword = false;
      }}>
      OK
    </button>
  </div>
{:else}
  <div class="container">
    <button
      class="download-btn"
      title="Descargar fichero torrent"
      on:click={downloadTorrent}>
      Descargar
    </button>
    <div class="buttons">
      <button
        title="Descargar en la carpeta por defecto"
        on:click={downloadTorrentOnFolder}>
        <Icon data={download} />
      </button>
      <button
        title="Copiar enlace del torrent al portapapeles"
        on:click={copyDownloadUrl}>
        <Icon data={copy} />
      </button>
      <button
        title="Copiar la contraseña para descomprimir al portapapeles"
        on:click={copyPassword}>
        <Icon data={key} />
      </button>
    </div>
  </div>
{/if}
