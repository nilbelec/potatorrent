$.postJSON = function (url, data) {
    return $.ajax({
        url: url,
        data: JSON.stringify(data),
        contentType: 'application/json',
        method: 'POST'
    });
}
$.delete = function (url) {
    return $.ajax({
        url: url,
        method: 'DELETE'
    });
}
$.success = function (text) {
    $('#notification-text').closest('.alert').removeClass('alert-danger alert-warning').addClass('alert-success');
    $('#notification-text').text(text);
    $('#notification').toast('dispose')
    $('#notification').toast('show');
}
$.error = function (text) {
    $('#notification-text').closest('.alert').removeClass('alert-success alert-warning').addClass('alert-danger');
    $('#notification-text').text(text);
    $('#notification').toast('dispose')
    $('#notification').toast('show');
}
$.warning = function (text) {
    $('#notification-text').closest('.alert').removeClass('alert-danger alert-success').addClass('alert-warning');
    $('#notification-text').text(text);
    $('#notification').toast('dispose')
    $('#notification').toast('show');
}

$.fn.tooltip.Constructor.Default.boundary = 'viewport';
$.fn.tooltip.Constructor.Default.delay = { "show": 500, "hide": 100 };