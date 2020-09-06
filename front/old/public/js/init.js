$.extend($.expr[':'], {
    'containsi': function(elem, i, match, array)
    {
      return (elem.textContent || elem.innerText || '').toLowerCase()
      .indexOf((match[3] || "").toLowerCase()) >= 0;
    }
  });

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

const toastParentTemplate = document.getElementById('toast-template-container');
const toastTemplate = toastParentTemplate.innerHTML;
toastParentTemplate.remove()

$.showNotification = function(text, type, icon) {
    const $toast = $(toastTemplate);
    $toast.find('.toast-body i.fa').addClass(icon);
    $toast.addClass('toast-' + type).find('.toast-body .toast-text').text(text);
    $('#notifications').prepend($toast);
    $toast.toast('show');
}

$.info = function (text) {
    $.showNotification(text, 'info', 'fa-info-circle');
}
$.success = function (text) {
    $.showNotification(text, 'success', 'fa-check-square-o');
}
$.error = function (text) {
    $.showNotification(text, 'danger', 'fa-exclamation-triangle');
}
$.warning = function (text) {
    $.showNotification(text, 'warning', 'fa-exclamation-circle');
}

$.fn.tooltip.Constructor.Default.boundary = 'viewport';
$.fn.tooltip.Constructor.Default.delay = { "show": 500, "hide": 100 };

$.fn.select2.defaults.set("theme", "bootstrap4");
$.fn.select2.defaults.set("width", "100%");
$.fn.select2.defaults.set("language", "es");