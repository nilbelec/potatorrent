export async function searchOptions() {
    const json = await fetchJSON("/options");
    return {
        categories: toValuesAndLabels(json.categorias),
        qualities: toValuesAndLabels(json.calidades),
    }
}

export async function searchSubCategories(category) {
    const json = await fetchJSON(`/subcategories?categoria=${category}`);
    return toValuesAndLabels(json);
}

export async function searchTorrents(cat, subCat, quality, q, page) {
    const url = `/search?categoria=${cat || ""}&subcategoria=${subCat || ""}&calidad=${quality || ""}&q=${q || ""}&pg=${page || ""}`;
    const json = await fetchJSON(url);
    if (json && json.data && json.data.torrents) {
        if (json.data.torrents["0"]) {
            json.data.torrents = Object.keys(json.data.torrents).flatMap(k => {
                return Object.keys(json.data.torrents[k]).map(j => {
                    let torrent = json.data.torrents[k][j];
                    torrent.name = extractName(torrent.torrentName);
                    torrent.season = extractSeason(torrent.torrentName);
                    torrent.firstEpisode = extractFirstEpisode(torrent.torrentName, torrent.season);
                    torrent.lastEpisode = extractLastEpisode(torrent.torrentName, torrent.season);
                    torrent.singleEpisode = torrent.lastEpisode == undefined;
                    torrent.imagen = torrent.imagen == '/pictures/c/thumbs/' ? '/pctn/library/content/template/images/no-imagen.jpg' : torrent.imagen;
                    return torrent;
                });
            });
        } else {
            json.data.torrents = [];
        }
    }
    return json;
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
    var myRegexp = /Cap.(\d+)/g;
    var match = myRegexp.exec(torrentName);
    return match && match[1] && match[1].replace(new RegExp(season), '') && Number(match[1].replace(new RegExp(season), ''));
}

function extractLastEpisode(torrentName, season) {
    if (isNaN(season))
        return;
    var myRegexp = /Cap.\d+_(\d+)/g;
    var match = myRegexp.exec(torrentName);
    return match && match[1] && match[1].replace(new RegExp(season), '') && Number(match[1].replace(new RegExp(season), ''));
}

async function fetchJSON(url) {
    const response = await fetch(url);
    return await response.json();
}

function toValuesAndLabels(obj) {
    return Object.keys(obj).map(key => ({
        value: key,
        label: obj[key]
    })).sort((a, b) => a.label.localeCompare(b.label));
}