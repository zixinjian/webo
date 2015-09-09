/**
 * Created by rick on 15/9/5.
 */
+function ($) {
    'use strict';
    var poplayer = function (element, options) {
        this.init('poplayer', element, options)
    }

    if (!$.fn.popover) throw new Error('poplayer requires bootstrap.js')

    poplayer.VERSION  = '0.0.1'

    poplayer.DEFAULTS = $.extend({}, $.fn.tooltip.Constructor.DEFAULTS, {
        placement: 'bottom',
        trigger: 'click',
        content: '',
        template: '<div class="popover" role="tooltip"><div class="arrow"></div><h3 class="popover-title"></h3><div class="popover-content"></div>' +
        '<div class="modal-footer"><button type="button" class="btn btn-default btn-sm" data-dismiss="modal">Close</button></div></div>'
    })


    // NOTE: POPOVER EXTENDS popover
    // ================================

    poplayer.prototype = $.extend({}, $.fn.popover.Constructor.prototype)

    poplayer.prototype.constructor = poplayer

    poplayer.prototype.getDefaults = function () {
        return poplayer.DEFAULTS
    }

    poplayer.prototype.setContent = function () {
        var $tip    = this.tip()
        var title   = this.getTitle()
        var content = this.getContent()

        $tip.find('.popover-title')[this.options.html ? 'html' : 'text'](title)
        if(this.options.url){
            $tip.find('.popover-content').children().detach().end().load(this.options.url)
        }else{
            $tip.find('.popover-content').children().detach().end()[ // we use append for html objects to maintain js events
                this.options.html ? (typeof content == 'string' ? 'html' : 'append') : 'text'
                ](content)

            $tip.removeClass('fade top bottom left right in')
        }

        // IE8 doesn't accept hiding via the `:empty` pseudo selector, we have to do
        // this manually by checking the contents.
        if (!$tip.find('.popover-title').html()) $tip.find('.popover-title').hide()
    }

    // POPOVER PLUGIN DEFINITION
    // =========================

    function Plugin(option) {
        return this.each(function () {
            var $this   = $(this)
            var data    = $this.data('bs.poplayer')
            var options = typeof option == 'object' && option

            if (!data && /destroy|hide/.test(option)) return
            if (!data) $this.data('bs.poplayer', (data = new poplayer(this, options)))
            if (typeof option == 'string') data[option]()
        })
    }

    var old = $.fn.poplayer

    $.fn.poplayer             = Plugin
    $.fn.poplayer.Constructor = poplayer


    // POPOVER NO CONFLICT
    // ===================

    $.fn.poplayer.noConflict = function () {
        $.fn.poplayer = old
        return this
    }

}(jQuery);