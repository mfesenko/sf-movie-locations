var markers = [];
var infowindow;
var map;
var locationToMovies = {};

function initMap() {
    var mapCenter = {lat: 37.7832508, lng: -122.450038};

    map = new google.maps.Map(document.getElementById('map'), {
        zoom: 13,
        center: mapCenter,
        mapTypeControl: false
    });

    infowindow = new(CustomInfoWindow())();
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
    infowindow.close();
    markers.forEach(hideMarker);
    markers = [];
    locationToMovies = [];
    if (data != null) {
        data.forEach(createMarkersForMovie);
    }
}

function createMarkersForMovie(movie) {
    if (movie.locations) {
        movie.locations.forEach(function (location) {
            var locationMovies = locationToMovies[location.description];
            if (typeof(locationMovies) == 'undefined') {
                locationMovies = [movie];
                locationToMovies[location.description] = locationMovies;

                var marker = createMarker(location);
                marker.addListener('click', function () {
                    var msgData = {
                        movies: locationToMovies[location.description],
                        location: location
                    }
                    infowindow.setContent(getInfoMessageForMovie(msgData));
                    infowindow.open(map, marker);
                });
                markers.push(marker)
            } else {
                locationMovies.push(movie);
            }
        });
    }
}

function createMarker(location) {
    return new google.maps.Marker({
        position: {
            lat: location.lat,
            lng: location.lng
        },
        map: map,
        title: movie.title,
        icon: {
            url: '/static/images/marker.png',
            size: new google.maps.Size(32, 32),
            scaledSize:new google.maps.Size(32, 32)
        }
    });
}

function getInfoMessageForMovie(msgData) {
    return $("#messageTemplate").tmpl(msgData).html();
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