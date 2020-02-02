$(function () {
    $.getJSON("/version")
        .done(function (response) {
            if (!response)
                return;
            $('#title').append('<span id="version" class="badge badge-danger">' + response.current + '</span>')
            if (response.current !== response.latest) {
                const msg = 'Hay una nueva versión disponible: ' + response.latest;
                $('#github-link')
                    .prop('href', 'https://github.com/nilbelec/potatorrent/releases')
                    .toggleClass('text-blink bg-success')
                    .prop('title', msg)
                    .tooltip('dispose').tooltip();
            }
        });

    $('#schedules').tooltip();
    $('#github-link').tooltip();

    $.getJSON("/options")
        .done(function (response) {
            if (!response)
                return;
            appendOptions('#categoria', response.categorias);
            appendOptions('#calidad', response.calidades);
            $('#search-input').prop('disabled', false);
            $('#categoria').prop('disabled', false).select2({
                placeholder: 'Elige una categoría',
                allowClear: true
            });
            $('#calidad').prop('disabled', false).select2({
                placeholder: 'Elige la calidad del contenido',
                allowClear: true
            });
        });

    $('#title').on('click', function () {
        $("html, body").animate({ scrollTop: 0 }, 500);
    })

    $('#subcategoria').select2({
        placeholder: 'Elige una subcategoría',
        allowClear: true
    });

    function appendOptions(selector, map) {
        if (!selector || !map)
            return;
        var sortable = [];
        for (var key in map) {
            sortable.push([key, map[key]]);
        }

        sortable.sort(function (a, b) {
            return a[1].localeCompare(b[1]);
        });
        for (var k in sortable) {
            $(selector).append('<option value="' + sortable[k][0] + '">' + sortable[k][1] + '</option>')
        }
    }

    $('#search-input').on('keyup', function (e) {
        if (e.which == 13)
            $('#search').trigger('click');
    });

    $('#categoria').on('change', function () {
        $('#subcategoria').prop('disabled', true)
        $('#subcategoria').html('<option selected value=""></option>');
        const categoria = $('#categoria').val();
        if (!categoria) {
            return;
        }
        $('#subcategoria').prop('disabled', false)
        $.getJSON("/subcategories", { categoria: categoria })
            .done(function (response) {
                if (!response)
                    return;
                appendOptions('#subcategoria', response);
            })
    })

    const torrentParentTemplate = document.getElementById('result-template-container');
    const torrentTemplate = torrentParentTemplate.innerHTML;
    torrentParentTemplate.remove()

    const loading = '<div class="loading-icon text-center"><i class="fa fa-spinner fa-pulse fa-3x fa-fw"></i></div>';

    $('#search').on('click', function () {
        $('#search-container').data('params', {
            categoria: $('#categoria').val(),
            categoriaTexto: $('#categoria option:selected').text(),
            subcategoria: $('#subcategoria').val(),
            subcategoriaTexto: $('#subcategoria option:selected').text(),
            calidad: $('#calidad').val(),
            calidadTexto: $('#calidad option:selected').text(),
            q: $('#search-input').val(),
            pg: 1,
        });
        $('#results-top-bar').hide();
        $('#search-results-info').empty();
        $('#search-results-pagination').stop(true).css('bottom', '-80px').empty();
        doSearch();
    });

    function doSearch() {
        $('#search').prop('disabled', true);
        $('#search-results-container').html(loading);
        var params = $('#search-container').data('params');
        $.getJSON("/search", params).done(function (response) {
            $('#search-results-container').empty()
            $('#results-top-bar').show();
            if (!response || !response.success) {
                $('#search-results-container').html('<div class="alert alert-danger mt-2"><i class="fa fa-exclamation-triangle fa-fw" aria-hidden="true"></i>Se ha producido un error...</div>')
                return;
            }
            if (!response.data.items) {
                $('#search-results-container').html('<div class="alert alert-info mt-2"><i class="fa fa-info-circle fa-fw" aria-hidden="true"></i>No se han encontrado resultados</div>')
                return;
            }
            const page = params.pg;
            const first = ((page - 1) * 30) + 1
            const last = first + response.data.items - 1;
            const total = response.data.total;
            $('#search-results-info').html('Mostrando resultados del <strong>' + first + ' al ' + last + '</strong> de un total de <strong>' + total + '</strong>');
            if (response.data.all > 1) {
                const ul = $('<ul class="pagination justify-content-center"></ul>');
                if (response.data.all > 9)
                    ul.append('<li class="page-item' + (page == 1 ? ' disabled' : '') + '"><a class="page-link" data-page="' + 1 + '" href="#"><i class="fa fa-angle-double-left" aria-hidden="true"></i></a></li>')
                ul.append('<li class="page-item' + (page == 1 ? ' disabled' : '') + '"><a class="page-link" data-page="' + (page - 1) + '" href="#"><i class="fa fa-angle-left" aria-hidden="true"></i></a></li>')
                for (var k = 0; k < 5; k++) {
                    var c = page - 2;
                    while (c < 1)
                        c++;
                    c += k;
                    if (c <= response.data.all)
                        ul.append('<li class="page-item' + (c == page ? ' active' : '') + '"><a class="page-link" data-page="' + c + '" href="#">' + c + '</a></li>')
                }
                ul.append('<li class="page-item' + (page == response.data.all ? ' disabled' : '') + '"><a class="page-link" data-page="' + (page + 1) + '" href="#"><i class="fa fa-angle-right" aria-hidden="true"></i></a></li>')
                if (response.data.all > 9)
                    ul.append('<li class="page-item' + (page == response.data.all ? ' disabled' : '') + '"><a class="page-link" data-page="' + response.data.all + '" href="#"><i class="fa fa-angle-double-right" aria-hidden="true"></i></a></li>')
                $('#search-results-pagination').html(ul).stop(true).animate({ bottom: "0" }, 600);
                $('#search-results-pagination').find('.page-link').on('click', function (e) {
                    e.preventDefault();
                    var params = $('#search-container').data('params');
                    params.pg = Number($(this).data('page'));
                    $('#search-container').data('params', params);
                    doSearch();
                })
            }
            
            for (let k in response.data.torrents["0"]) {
                const t = response.data.torrents["0"][k];

                const $cont = $(torrentTemplate);
                $cont.find('.torrent-img').attr('data-src', '/image?path=' + (t.imagen == '/pictures/f/thumbs/' ? '/d20/library/content/template/images/no-imagen.jpg' : t.imagen))
                $cont.find('.torrent-title').text(t.torrentName)
                $cont.find('.torrent-calidad').text(t.calidad)
                $cont.find('.torrent-fecha').text(t.torrentDateAdded)
                $cont.find('.torrent-tamano').text(t.torrentSize)
                $cont.find('.download-group').data('id', t.torrentID).data('guid', t.guid).data('date', t.torrentDateAdded)

                $('#search-results-container').append($cont)

                $cont.find('.has-tooltip').tooltip()

                $cont.find('.torrent-download, .torrent-download-copy, .torrent-download-password').on('click', function (e) {
                    e.stopPropagation()
                    const $btn = $(this);
                    const $group = $btn.closest('.download-group');
                    if ($group.data('download-data')) {
                        processDownloadData($btn);
                        return;
                    }
                    $btn.prop('disabled', true);
                    const innerHTML = $btn.html();
                    $btn.html('<span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span> Buscando...');

                    const id = $group.data('id');
                    const guid = $group.data('guid');
                    const date = $group.data('date');

                    $.getJSON('/download', { id: id, guid: guid, date: date })
                        .done(function (response) {
                            $group.data('download-data', response);
                            processDownloadData($btn);
                        }).always(function () {
                            $btn.prop('disabled', false);
                            $btn.html(innerHTML);
                        })
                })

                function processDownloadData($btn) {
                    const $group = $btn.closest('.download-group');
                    const data = $group.data('download-data')
                    if (!data || !data.url) {
                        $.warning('Parece que el torrent ya no existe...');
                        return;
                    }
                    if ($btn.hasClass('torrent-download')) {
                        window.open(data.url);
                        if (data.password) {
                            $('#password').text(data.password);
                            $('#password-modal').modal('show');
                        }
                    } else if ($btn.hasClass('torrent-download-copy')) {
                        copyToClipboard(data.url)
                        $.success('Enlace copiado al portapapeles :)');
                    } else if ($btn.hasClass('torrent-download-password')) {
                        if (!data.password) {
                            $.info('Este torrent no parece necesitar ninguna contraseña');
                            return;
                        }
                        $('#password').text(data.password);
                        $('#password-modal').modal('show');
                    }
                }

                const copyToClipboard = str => {
                    const el = document.createElement('textarea');
                    el.value = str;
                    el.setAttribute('readonly', '');
                    el.style.position = 'absolute';
                    el.style.left = '-9999px';
                    document.body.appendChild(el);
                    el.select();
                    document.execCommand('copy');
                    document.body.removeChild(el);
                };
            }
            lazyload();
        }).fail(function(){
            $('#search-results-container').empty()
            $.error("Ups! Se ha producido un error...");
        }).always(function (){
            $('#search').prop('disabled', false);
        })
    }

    $('#schedule-search').on('click', function () {
        $('#schedule-search-modal').modal('show');
    })

    $('#schedule-add-folder').select2({
        dropdownParent: $("#schedule-search-modal"),
        placeholder: 'Ruta del directorio',
        ajax: {
            url: '/folder',
            dataType: 'json',
            delay: 100
        }
    });

    $('#schedule-search-modal').on('show.bs.modal', function () {
        const params = $('#search-container').data('params');
        $('#staticCategoria').text(params.categoriaTexto ? params.categoriaTexto : "-");
        $('#staticSubcategoria').text(params.subcategoriaTexto ? params.subcategoriaTexto : "-");
        $('#staticCalidad').text(params.calidadTexto ? params.calidadTexto : "-");
        $('#staticPalabras').text(params.q ? params.q : "-");

        $('#schedule-add-folder').removeClass('is-invalid');
        $('#schedule-add-interval').removeClass('is-invalid');
        $('#schedule-add-error').hide();
    })
    $('#schedule-search-modal').on('shown.bs.modal', function () {
        $('#schedule-add-folder').select()
    });

    $('#schedule-search-save').on('click', function () {
        $('#schedule-add-folder').removeClass('is-invalid');
        $('#schedule-add-interval').removeClass('is-invalid');
        $('#schedule-add-error').hide();

        const params = $('#search-container').data('params');

        let ok = true;
        const folder = $('#schedule-add-folder').val();
        if (!folder) {
            $('#schedule-add-folder').addClass('is-invalid');
            ok = false;
        }
        const interval = $('#schedule-add-interval').val();
        if (!isPositiveInteger(interval)) {
            $('#schedule-add-interval').addClass('is-invalid');
            ok = false;
        }
        if (!ok)
            return;

        const $btn = $(this);
        $btn.prop('disabled', true);
        const innerHTML = $btn.html();
        $btn.html('<span class="spinner-border spinner-border-sm"></span> Guardando...');

        const url = '/schedule';
        $.postJSON(url, {
            folder: folder,
            interval: Math.floor(Number(interval)),
            params: params,
        }).done(function (response) {
            refreshSchedules();
            $('#schedule-search-modal').modal('hide');
        }).fail(function (xhr) {
            if (!xhr || !xhr.responseText)
                $('#schedule-add-error').html('<i class="fa fa-exclamation-triangle fa-fw" aria-hidden="true"></i> Ups! Se ha producido un error...').show();
            else
                $('#schedule-add-error').html('<i class="fa fa-exclamation-triangle fa-fw" aria-hidden="true"></i> ' + xhr.responseText).show();
        }).always(function () {
            $btn.prop('disabled', false).html(innerHTML);
        });
    });


    const schedulesParentTemplate = document.getElementById('schedule-template-container');
    const scheduleTemplate = schedulesParentTemplate.innerHTML;
    schedulesParentTemplate.remove()

    let req;
    let to;

    $('#schedules').on('click', function () {
        $('#schedules-modal').modal('show');
    });

    $('#schedules-refresh-btn').on('click', function () {
        refreshSchedules(true);
    });

    $('#schedules-modal').on('show.bs.modal', function(){
        refreshSchedules(true);
    })

    $('#schedules-modal').on('hidden.bs.modal', function(){
        to = setTimeout(refreshSchedules, schedulesRefreshTime);
    })

    const schedulesRefreshTime = 10000;
    function refreshSchedules(show) {
        if (to) {
            clearTimeout(to);
        }
        if (req) {
            req.abort();
        }
        const url = '/schedules';
        $('#schedules-refresh-btn i').addClass('fa-spinner');
        req = $.getJSON(url).done(function (response) {
            $('#schedules-results').empty();
            $('#schedules-refresh-btn i').removeClass('fa-spinner');

            if (!response || !response.length) {
                $('#schedules-count').hide();
                $('#schedules-none').show();
                return;
            }
            $('#schedules-none').hide();
            $('#schedules-count').text(response.length).show();
            if (!show) {
                to = setTimeout(refreshSchedules, schedulesRefreshTime);
                return;
            }
            const $tbody = $('#schedules-table tbody').empty();
            for (var k in response) {
                const s = response[k];

                const $cont = $(scheduleTemplate);
                if (s.params.categoriaTexto)
                    $cont.find('.schedule-categoria').text(s.params.categoriaTexto);
                else
                    $cont.find('.schedule-categoria').hide();
                if (s.params.subcategoriaTexto)
                    $cont.find('.schedule-subcategoria').text(s.params.subcategoriaTexto);
                else
                    $cont.find('.schedule-subcategoria').hide();
                if (s.params.calidadTexto)
                    $cont.find('.schedule-calidad').text(s.params.calidadTexto);
                else
                    $cont.find('.schedule-calidad').hide();
                if (s.params.q)
                    $cont.find('.schedule-palabras').text(s.params.q);
                else
                    $cont.find('.schedule-palabras').hide();

                $cont.find('.schedule-folder').prop('title', 'Carpeta de descarga: ' + s.folder).tooltip();
                $cont.find('.schedule-interval').text(s.interval);
                $cont.find('.schedule-last-execution').text(moment(s.lastExecutionTime).fromNow());

                if (!s.lastTorrentImage || s.lastTorrentImage == '/pictures/f/thumbs/') {
                    $cont.find('.schedule-img').attr('data-src', '/image?path=/pctn/library/content/template/images/no-imagen.jpg');
                } else {
                    $cont.find('.schedule-img').attr('data-src', '/image?path=' + s.lastTorrentImage)
                }
                
                
                if (s.lastTorrentName) {
                    $cont.find('.schedule-last-torrent').text(s.lastTorrentName);
                    $cont.find('.schedule-last-torrent-date').text(s.lastTorrentDate);
                }

                if (s.disabled) {
                    $cont.addClass('schedule-disabled');
                    $cont.find('.schedule-disable').hide();
                } else {
                    $cont.find('.schedule-enable').hide();
                }

                $('#schedules-results').append($cont);
                $cont.find('.has-tooltip').tooltip();

                $cont.data('schedule-id', s.id);
                $cont.find('.schedule-enable').on('click', function () {
                    const $cont = $(this).closest('.schedule');
                    const id = $cont.data('schedule-id');
                    const $btn = $(this);
                    const innerHTML = $btn.html();
                    $btn.prop('disabled', true).html('<span class="spinner-border spinner-border-sm"></span>');
                    const url = '/schedule/enable?id=' + id;
                    $.post(url)
                        .done(function () {
                            $cont.removeClass('schedule-disabled');
                            $btn.hide().tooltip('dispose');
                            $cont.find('.schedule-disable').show().tooltip();
                        }).fail(function () {
                            $.error("Ups! Se ha producido un error...")
                        }).always(function () {
                            $btn.prop('disabled', false).html(innerHTML);
                        })

                });

                $cont.find('.schedule-disable').on('click', function () {
                    const $cont = $(this).closest('.schedule');
                    const id = $cont.data('schedule-id');
                    const $btn = $(this);
                    $btn.prop('disabled', true);
                    const innerHTML = $btn.html();
                    $btn.html('<span class="spinner-border spinner-border-sm"></span>');
                    const url = '/schedule/disable?id=' + id;
                    $.post(url)
                        .done(function () {
                            $cont.addClass('schedule-disabled');
                            $btn.hide().tooltip('dispose');
                            $cont.find('.schedule-enable').show().tooltip();
                        }).fail(function () {

                        }).always(function () {
                            $btn.prop('disabled', false).html(innerHTML);
                        })
                });
                $cont.find('.schedule-delete').on('click', function () {
                    const $cont = $(this).closest('.schedule');
                    const id = $cont.data('schedule-id');
                    const $btn = $(this);
                    $btn.prop('disabled', true);
                    const innerHTML = $btn.html();
                    $btn.html('<span class="spinner-border spinner-border-sm"></span>');
                    const url = '/schedule?id=' + id;
                    $.delete(url)
                        .done(function () {
                            $btn.tooltip('dispose');
                            refreshSchedules(true);
                        }).fail(function () {
                            $.error("Ups! Se ha producido un error...")
                        }).always(function () {
                            $btn.prop('disabled', false).html(innerHTML);
                        })
                });
            }
            $('#schedules-results').show();
            lazyload();
            filterSchedules();
        });
    }
    refreshSchedules();

    $('#schedules-filter').on('keyup', function(){
        filterSchedules();
    })
    function filterSchedules(){
        $('#schedules-modal .schedule').show();
        const filter = $('#schedules-filter').val().trim();
        if (!filter)
            return;
        $('#schedules-modal .schedule:not(:containsi(' + filter + '))').hide();
    }

    function isPositiveInteger(str) {
        var n = Math.floor(Number(str));
        return n !== Infinity && String(n) === str && n >= 0;
    }

    $('#search').trigger('click');
});