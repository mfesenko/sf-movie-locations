window.SfMovies = new (function () {
    var exports = {};
    var markers = [];
    var infowindow;
    var map;
    var locationToMovies = {};

    exports.initMap = function () {
        var mapCenter = {lat: 37.7832508, lng: -122.450038};

        map = new google.maps.Map(document.getElementById('map'), {
            zoom: 13,
            center: mapCenter,
            mapTypeControl: false
        });

        infowindow = new (CustomInfoWindow())();
        google.maps.event.addListenerOnce(map, 'idle', function () {
            initEventHandlers();
            loadAllMovieLocations();
        });
    }

    exports.joinWithComma = function(values) {
        return values.join(', ');
    }

    function initEventHandlers() {
        initAutocomplete();
        $('#filter').click(filter);
        $('#movie').keyup(function (event) {
            var keyCode = event.keyCode ? event.keyCode : event.which;
            if (13 === keyCode) {
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
            },
            onSelect: function () {
                filter();
            }
        });
    }

    function onMovieLocationsLoaded(data) {
        infowindow.close();
        resetMarkers();
        data && data.forEach(createMarkersForMovie);
    }

    function resetMarkers() {
        markers.forEach(hideMarker);
        markers = [];
        locationToMovies = [];
    }

    function createMarkersForMovie(movie) {
        if (!movie.locations) return;

        movie.locations.forEach(function (location) {
            var locationMovies = locationToMovies[location.description];

            if (locationMovies) {
                locationMovies.push(movie);
            } else {
                locationToMovies[location.description] = [movie];

                var marker = createMarker(location);
                marker.addListener('click', function () {
                    var msgData = {
                        movies: locationToMovies[location.description],
                        location: location
                    };
                    infowindow.setContent(getInfoMessageForMovie(msgData));
                    infowindow.open(map, marker);
                });
                markers.push(marker)
            }
        });
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
                scaledSize: new google.maps.Size(32, 32)
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

    return exports;
});