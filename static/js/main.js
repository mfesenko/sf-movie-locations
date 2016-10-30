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
    $('#filter').click(filter);
    $('#movie').keyup(function (e) {
        if (e.keyCode === 13) {
            filter();
        }
    })
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
                        value: movie.title,
                        data: movie.id
                    };
                })
            };
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
                message = getInfoMessageForMovie(movie);
                infowindow.setContent(message);
                infowindow.open(map, marker);
            });
            markers.push(marker)
        });
    }
}

function getInfoMessageForMovie(movie) {
    return $("#messageTemplate").tmpl(movie).html();
}

function filter() {
    $('#filter').blur();
    filterMovieLocations($('#movie').val());
}

function loadAllMovieLocations() {
    $.get('/api/movieLocations',
        onMovieLocationsLoaded,
        'json'
    );
}

function filterMovieLocations(movieTitle) {
    $.get('/api/movieLocations?title=' + movieTitle,
        onMovieLocationsLoaded,
        'json'
    );
}

function hideMarker(marker) {
    marker.setMap(null);
}

function joinWithComma(values) {
    return values.join(', ');
}