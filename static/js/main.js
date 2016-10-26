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


    $.get('/api/movieLocations',
        onMovieLocationsLoaded,
        'json'
    );
}


function onMovieLocationsLoaded(data) {
    data.forEach(createMarkersForMovie);
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
                message = 'Movie: ' + movie.title+ ' (' + movie.year + ')';
                infowindow.setContent(message);
                infowindow.open(map, marker);
            });
            markers.push(marker)
        });
    }
}
