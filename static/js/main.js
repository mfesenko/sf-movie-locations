var markers = [];
var infowindow;
var map;

function initMap() {
    var mapCenter = {lat: 37.7832508, lng: -122.450038};

    map = new google.maps.Map(document.getElementById('map'), {
        zoom: 13,
        center: mapCenter,
        mapTypeControl: false
    });

    infowindow = new google.maps.InfoWindow({
        content: ''
    });
    google.maps.event.addListenerOnce(map, 'idle', function () {
        initEventHandlers();
        loadAllMovieLocations();
    });
}

function initEventHandlers() {
    initAutocomplete();
    $('#reset').click(reset);
}

function initAutocomplete() {
    $('#movie').autocomplete({
        serviceUrl: '/api/movies/',
        paramName: 'title',
        transformResult: function (response) {
            var movies = JSON.parse(response);
            return {
                suggestions: $.map(movies, function (movie) {
                    return {
                        value: movie.title + " (" + movie.year + ")",
                        data: movie.id
                    };
                })
            };
        },
        onSelect: function (suggestion) {
            filterMovieLocations(suggestion.data);
        }
    });
}

function onMovieLocationsLoaded(data) {
    markers.forEach(hideMarker);
    markers = [];
    if (data != null) {
        data.forEach(createMarkersForMovie);
    }
}

function createMarkersForMovie(movie) {
    if (movie.locations) {
        movie.locations.forEach(function (location) {
            var marker = new google.maps.Marker({
                position: location,
                map: map,
                title: movie.title,
                icon: '/static/images/marker.png'
            });
            marker.addListener('click', function () {
                message = 'Movie: ' + movie.title + ' (' + movie.year + ')';
                infowindow.setContent(message);
                infowindow.open(map, marker);
            });
            markers.push(marker)
        });
    }
}

function reset() {
    $('#reset').blur();
    $('#movie').val('');
    loadAllMovieLocations();
}

function loadAllMovieLocations() {
    $.get('/api/movieLocations',
        onMovieLocationsLoaded,
        'json'
    );
}

function filterMovieLocations(movieId) {
    $.get('/api/movieLocations/' + movieId,
        onMovieLocationLoaded,
        'json'
    );
}

function onMovieLocationLoaded(data) {
    markers.forEach(hideMarker);
    markers = [];
    createMarkersForMovie(data);
}

function hideMarker(marker) {
    marker.setMap(null);
}