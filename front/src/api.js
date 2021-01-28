export async function getVersion() {
    return await fetchJSON("/api/version");
}

export async function isNewVersion() {
    const latest = await getLatestVersion();
    const current = await getVersion();
    return latest &&
        current &&
        latest.version !== current.version &&
        latest.version;
}

export async function getLatestVersion() {
    return await fetchJSON("/api/latest");
}

export async function getConfiguration() {
    return await fetchJSON("/api/configuration");
}

export async function saveConfiguration(conf) {
    const resp = await fetch('/api/configuration', {
        method: 'POST',
        body: JSON.stringify(conf)
    });
    if (!resp.ok) {
        const text = await resp.text();
        throw new Error(text)
    }
}

export async function getSchedules() {
    const json = await fetchJSON("/api/schedules");
    if (!json || !json.length)
        return json;
    return json.map(s => {
        let lastExecutionTime = new Date(s.lastExecutionTime);
        const sch = {
            id: s.id,
            params: s.params,
            lastExecutionTime: lastExecutionTime.getFullYear() > 1 ? timeAgo(lastExecutionTime) : "",
            error: s.error,
            disabled: s.disabled
        }
        if (!s.lastTorrentName)
            return sch;
        const season = extractSeason(s.lastTorrentName);
        const firstEpisode = extractFirstEpisode(s.lastTorrentName, season);
        const lastEpisode = extractLastEpisode(s.lastTorrentName, season);
        sch.lastTorrent = {
            id: s.lastTorrentID,
            name: extractName(s.lastTorrentName),
            season: season,
            firstEpisode: firstEpisode,
            lastEpisode: lastEpisode,
            singleEpisode: lastEpisode == undefined,
            imagen: prepareImage(s.lastTorrentImage),
            date: s.lastTorrentDate,
            password: s.lastTorrentPassword,
            fullName: s.lastTorrentName,
        }
        return sch;
    })
}

function timeAgo(date) {
    const now = new Date();
    let secs = Math.floor((now.getTime() - date.getTime()) / 1000);
    if (secs < 30)
        return "Hace unos segundos";
    if (secs < 60)
        return `Hace medio minuto`
    if (secs < 120)
        return `Hace un minuto`
    const mins = Math.floor(secs / 60);
    if (mins < 60)
        return `Hace ${mins} minutos`
    if (mins < 120)
        return `Hace una hora`
    const hours = Math.floor(mins / 60);
    if (hours < 24)
        return `Hace ${hours} horas`
    if (hours < 48)
        return `Hace un día`
    const days = Math.floor(hours / 24);
    return `Hace ${days} días`
}

export async function enableSchedule(id) {
    await fetch('/api/schedule/enable?id=' + id, {
        method: 'POST'
    });
}

export async function disableSchedule(id) {
    await fetch('/api/schedule/disable?id=' + id, {
        method: 'POST'
    });
}

export async function deleteSchedule(id) {
    await fetch('/api/schedule?id=' + id, {
        method: 'DELETE'
    });
}

export async function saveSearch(searchParams) {
    const data = {
        "params": {
            "categoria": searchParams && searchParams.category && searchParams.category.value,
            "categoriaTexto": searchParams && searchParams.category && searchParams.category.label,
            "subcategoria": searchParams && searchParams.subCategory && searchParams.subCategory.value,
            "subcategoriaTexto": searchParams && searchParams.subCategory && searchParams.subCategory.label,
            "calidad": searchParams && searchParams.quality && searchParams.quality.value,
            "calidadTexto": searchParams && searchParams.quality && searchParams.quality.label,
            "q": searchParams && searchParams.text
        }
    }
    return await fetch('/api/schedule', {
        method: 'POST',
        body: JSON.stringify(data)
    });
}

export async function searchOptions() {
    const json = await fetchJSON("/api/options");
    return {
        categories: toValuesAndLabels(json.categorias),
        qualities: toValuesAndLabels(json.calidades),
    }
}

export async function searchSubCategories(category) {
    const json = await fetchJSON(`/api/subcategories?categoria=${category}`);
    return toValuesAndLabels(json);
}

export async function searchTorrents(cat, subCat, quality, q, page) {
    const url = `/api/search?categoria=${cat || ""}&subcategoria=${subCat || ""}&calidad=${quality || ""}&q=${q || ""}&pg=${page || ""}`;
    const json = await fetchJSON(url);
    if (json && json.data && json.data.torrents) {
        if (json.data.torrents["0"]) {
            json.data.torrents = Object.keys(json.data.torrents).flatMap(k => {
                return Object.keys(json.data.torrents[k]).map(j => {
                    let torrent = json.data.torrents[k][j];
                    torrent.fullName = torrent.torrentName;
                    torrent.date = torrent.torrentDateAdded;
                    torrent.name = extractName(torrent.torrentName);
                    torrent.season = extractSeason(torrent.torrentName);
                    torrent.firstEpisode = extractFirstEpisode(torrent.torrentName, torrent.season);
                    torrent.lastEpisode = extractLastEpisode(torrent.torrentName, torrent.season);
                    torrent.singleEpisode = torrent.lastEpisode == undefined;
                    torrent.imagen = prepareImage(torrent.imagen);
                    return torrent;
                });
            });
        } else {
            json.data.torrents = [];
        }
    }
    return json;
}

function getTorrentDownloadQueryParams(torrent) {
    let url = `id=${torrent.torrentID}&guid=${torrent.guid}&date=${torrent.torrentDateAdded}`;
    if (torrent.season)
        url += `&season=${torrent.season}`
    if (torrent.firstEpisode)
        url += `&firstEpisode=${torrent.firstEpisode}`
    if (torrent.lastEpisode)
        url += `&lastEpisode=${torrent.lastEpisode}`
    return url;
}

export function getDownloadInfo(torrent) {
    const params = getTorrentDownloadQueryParams(torrent);
    return fetchJSON(`/api/download/info?${params}`);
}

export function getDownloadLink(url) {
    const urlEnc = encodeURIComponent(url);
    return `/api/download/file?url=${urlEnc}`;
}

export function downloadOnFolder(torrent) {
    const params = getTorrentDownloadQueryParams(torrent);
    return fetchJSON(`/api/download/onFolder?${params}`);
}

function prepareImage(img) {
    if (!img || img == '/pictures/c/thumbs/')
        return '/pctn/library/content/template/images/no-imagen.jpg';
    return img;
}

function extractName(torrentName) {
    return torrentName.replace(/\[.+/, '').replace(/\(\d\d\d\d\)/, '').replace(/ \- Temporada.+/, '').trim();
}

function extractSeason(torrentName) {
    var myRegexp = /Temporada (\d+)/g;
    var match = myRegexp.exec(torrentName);
    return match && match[1] && Number(match[1]);
}

function extractFirstEpisode(torrentName, season) {
    if (isNaN(season))
        return;
    var myRegexp = /Cap\.\s?(\d+)/g;
    var match = myRegexp.exec(torrentName);
    return match && match[1] && match[1].replace(new RegExp(season), '') && Number(match[1].replace(new RegExp(season), ''));
}

function extractLastEpisode(torrentName, season) {
    if (isNaN(season))
        return;
    var myRegexp = /Cap\.\s?\d+_(\d+)/g;
    var match = myRegexp.exec(torrentName);
    return match && match[1] && match[1].replace(new RegExp(season), '') && Number(match[1].replace(new RegExp(season), ''));
}

async function fetchJSON(url) {
    const resp = await fetch(url);
    if (!resp.ok) {
        const text = await resp.text();
        throw new Error(text)
    }
    return await resp.json();
}

function toValuesAndLabels(obj) {
    return Object.keys(obj).map(key => ({
        value: key,
        label: obj[key]
    })).sort((a, b) => a.label.localeCompare(b.label));
}