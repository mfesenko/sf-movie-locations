/**
 * This code was borrowed from https://gist.github.com/SaneMethod/9140722#file-customwindow-jquery-js
 *
 * Create a custom overlay for our window marker display, extending google.maps.OverlayView.
 * This is somewhat complicated by needing to async load the google.maps api first - thus, we
 * wrap CustomWindow into a closure, and when instantiating CustomWindow, we first execute the closure (to create
 * our CustomWindow function, now properly extending the newly loaded google.maps.OverlayView), and then
 * instantiate said function.
 * Note that this version uses jQuery.
 * @type {Function}
 */
var CustomInfoWindow = function() {
    var CustomWindow = function(){
        this.container = $('<div class="map-info-window"></div>');
        this.layer = null;
        this.marker = null;
        this.position = null;
    };
    /**
     * Inherit from OverlayView
     * @type {google.maps.OverlayView}
     */
    CustomWindow.prototype = new google.maps.OverlayView();
    /**
     * Called when this overlay is set to a map via this.setMap. Get the appropriate map pane
     * to add the window to, append the container, bind to close element.
     * @see CustomWindow.open
     */
    CustomWindow.prototype.onAdd = function(){
        this.layer = $(this.getPanes().floatPane);
        this.layer.append(this.container);
        var window = this;
        this.container.find('.map-info-close').on('click', function(event){
            // Close info window on click
            window.close();
            if (event) {
                event.stopPropagation();
            }
        });
    };
    /**
     * Called after onAdd, and every time the map is moved, zoomed, or anything else that
     * would effect positions, to redraw this overlay.
     */
    CustomWindow.prototype.draw = function(){
        var markerIcon = this.marker.getIcon(),
            cHeight = this.container.outerHeight() + markerIcon.scaledSize.height + 10,
            cWidth = this.container.width() / 2 + markerIcon.scaledSize.width / 2;
        this.position = this.getProjection().fromLatLngToDivPixel(this.marker.getPosition());
        this.container.css({
            'top':this.position.y - cHeight,
            'left':this.position.x - cWidth
        });
    };
    /**
     * Called when this overlay has its map set to null.
     * @see CustomWindow.close
     */
    CustomWindow.prototype.onRemove = function(){
        this.container.remove();
    };
    /**
     * Sets the contents of this overlay.
     * @param {string} html
     */
    CustomWindow.prototype.setContent = function(html){
        this.container.html(html);
    };
    /**
     * Sets the map and relevant marker for this overlay.
     * @param {google.maps.Map} map
     * @param {google.maps.Marker} marker
     */
    CustomWindow.prototype.open = function(map, marker){
        this.marker = marker;
        this.setMap(map);
    };
    /**
     * Close this overlay by setting its map to null.
     */
    CustomWindow.prototype.close = function(){
        this.setMap(null);
    };
    return CustomWindow;
};