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
    addCustomControlsToMap(map);
    google.maps.event.addListenerOnce(map, 'idle', function () {
        initAutocomplete();
        loadAllMovieLocations();
    });
}

function addCustomControlsToMap(map) {
    var filterControlDiv = buildFilterControlDiv();
    filterControlDiv.index = 1;
    map.controls[google.maps.ControlPosition.TOP_RIGHT].push(filterControlDiv);
}

function buildFilterControlDiv() {
    var filterControlDiv = document.createElement('div');
    filterControlDiv.index = 1;
    filterControlDiv.appendChild(buildInputGroupDiv());
    return filterControlDiv;
}

function buildInputGroupDiv() {
    var inputGroupDiv = document.createElement('div');
    inputGroupDiv.className = 'filter input-group';
    inputGroupDiv.appendChild(buildFilterInput());
    inputGroupDiv.appendChild(buildInputGroupSpan());
    return inputGroupDiv;
}

function buildFilterInput() {
    var controlInput = document.createElement('input');
    controlInput.type = "text";
    controlInput.className = 'form-control';
    controlInput.id = 'movie';
    controlInput.placeholder = 'Filter by movie name...';
    return controlInput;
}

function buildInputGroupSpan() {
    var inputGroupSpan = document.createElement('span');
    inputGroupSpan.className = 'input-group-btn';
    inputGroupSpan.appendChild(buildResetButton());
    return inputGroupSpan;
}

function buildResetButton() {
    var resetButton = document.createElement('button');
    resetButton.id = 'reset';
    resetButton.innerText = 'Reset';
    resetButton.className = 'btn btn-default';
    resetButton.addEventListener('click', reset);
    return resetButton;
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
                title: movie.title
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